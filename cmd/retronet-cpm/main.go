// Comando retronet-cpm: shell e runner CP/M-like sopra retronet-8080.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/conformance"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
	"github.com/retronet-labs/retronet-cpm/shell"
)

type runConfig struct {
	diskPath    string
	runPath     string
	steps       uint64
	trace       bool
	traceJSON   string
	input       string
	aluName     string
	writeDisk   bool
	conformance bool
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(stderr, "errore: %v\n", err)
		return 2
	}
	alu, err := parseALU(cfg.aluName)
	if err != nil {
		fmt.Fprintf(stderr, "errore ALU: %v\n", err)
		return 2
	}
	traceSink, closeTrace, err := configureTrace(cfg, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "errore trace: %v\n", err)
		return 1
	}
	if closeTrace != nil {
		defer func() {
			if err := closeTrace(); err != nil {
				fmt.Fprintf(stderr, "errore chiusura trace: %v\n", err)
			}
		}()
	}
	if cfg.conformance {
		return runConformance(stdout)
	}
	drive, err := openDrive(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "errore drive A: %v\n", err)
		return 2
	}
	input := stdin
	if cfg.input != "" {
		input = strings.NewReader(cfg.input)
	}
	if cfg.runPath != "" {
		return runProgram(cfg, drive, alu, traceSink, input, stdout, stderr)
	}
	sh, err := shell.New(shell.Config{
		Drive:     drive,
		Input:     input,
		Output:    stdout,
		ALU:       alu,
		StepLimit: cfg.steps,
		Trace:     traceSink,
	})
	if err != nil {
		fmt.Fprintf(stderr, "errore shell: %v\n", err)
		return 2
	}
	if err := sh.Run(); err != nil {
		fmt.Fprintf(stderr, "errore shell: %v\n", err)
		return 1
	}
	return 0
}

func parseFlags(args []string, stderr io.Writer) (runConfig, error) {
	fs := flag.NewFlagSet("retronet-cpm", flag.ContinueOnError)
	fs.SetOutput(stderr)
	cfg := runConfig{
		diskPath: ".",
		steps:    cpm.DefaultStepLimit,
		aluName:  defaultALUName(),
	}
	fs.StringVar(&cfg.diskPath, "disk", cfg.diskPath, "directory host mappata come drive A:")
	fs.StringVar(&cfg.runPath, "run", "", "programma .COM da eseguire e poi terminare")
	fs.Uint64Var(&cfg.steps, "steps", cfg.steps, "limite massimo di istruzioni 8080")
	fs.BoolVar(&cfg.trace, "trace", false, "stampa trace testuale delle istruzioni e chiamate BDOS")
	fs.StringVar(&cfg.traceJSON, "trace-json", "", "scrive trace JSON Lines nel file indicato")
	fs.StringVar(&cfg.input, "input", "", "input testuale usato da programma o shell")
	fs.StringVar(&cfg.aluName, "alu", cfg.aluName, "backend ALU: native o gate")
	fs.BoolVar(&cfg.writeDisk, "write-disk", false, "abilita funzioni BDOS che modificano il drive host")
	fs.BoolVar(&cfg.conformance, "conformance", false, "esegue la suite sintetica integrata")
	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func openDrive(cfg runConfig) (disk.Drive, error) {
	if cfg.writeDisk {
		return disk.NewWritableHostDrive(cfg.diskPath)
	}
	return disk.NewHostDrive(cfg.diskPath)
}

func defaultALUName() string {
	if value := strings.TrimSpace(os.Getenv("RETRONET_CPM_ALU")); value != "" {
		return value
	}
	return "native"
}

func parseALU(name string) (cpu.ALUBackend, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "native", "":
		return cpu.Native, nil
	case "gate":
		return cpu.Gate, nil
	default:
		return nil, fmt.Errorf("valore %q non valido, usa native o gate", name)
	}
}

func runProgram(cfg runConfig, drive disk.Drive, alu cpu.ALUBackend, trace cpm.TraceSink, input io.Reader, stdout io.Writer, stderr io.Writer) int {
	name, data, err := loadProgram(drive, cfg.runPath)
	if err != nil {
		fmt.Fprintf(stderr, "errore caricamento programma: %v\n", err)
		return 1
	}
	console := bdos.NewStreamConsole(input, stdout)
	m, err := cpm.NewMachine(cpm.Config{
		ALU:       alu,
		StepLimit: cfg.steps,
		Console:   console,
		Disk:      drive,
		Trace:     trace,
	})
	if err != nil {
		fmt.Fprintf(stderr, "errore macchina: %v\n", err)
		return 2
	}
	if err := m.LoadCOM(name, data); err != nil {
		fmt.Fprintf(stderr, "errore caricamento COM: %v\n", err)
		return 1
	}
	result, err := m.Run(cfg.steps)
	printDump(stdout, name, result, m.CPU)
	if err != nil {
		fmt.Fprintf(stderr, "errore esecuzione: %v\n", err)
		return 1
	}
	return 0
}

func loadProgram(drive disk.Drive, name string) (string, []byte, error) {
	if hasPathSeparator(name) || filepath.IsAbs(name) {
		data, err := os.ReadFile(name)
		return filepath.Base(name), data, err
	}
	if data, err := os.ReadFile(name); err == nil {
		return filepath.Base(name), data, nil
	}
	driveName := name
	if !strings.Contains(driveName, ".") {
		driveName += ".COM"
	}
	data, err := drive.ReadFile(driveName)
	return driveName, data, err
}

func hasPathSeparator(value string) bool {
	return strings.ContainsAny(value, `/\`)
}

func configureTrace(cfg runConfig, stdout io.Writer) (cpm.TraceSink, func() error, error) {
	var sinks []cpm.TraceSink
	if cfg.trace {
		sinks = append(sinks, func(event cpm.TraceEvent) {
			if event.Kind == cpm.TraceBDOS {
				fmt.Fprintf(stdout, "trace=%d bdos C=%d pc=0x%04X\n", event.Sequence, event.BDOSFunction, event.PC)
				return
			}
			fmt.Fprintf(stdout, "trace=%d %s\n", event.Sequence, event.Disassembly)
		})
	}
	var closeTrace func() error
	if cfg.traceJSON != "" {
		file, err := os.Create(cfg.traceJSON)
		if err != nil {
			return nil, nil, err
		}
		encoder := json.NewEncoder(file)
		sinks = append(sinks, func(event cpm.TraceEvent) {
			_ = encoder.Encode(event)
		})
		closeTrace = file.Close
	}
	if len(sinks) == 0 {
		return nil, closeTrace, nil
	}
	return func(event cpm.TraceEvent) {
		for _, sink := range sinks {
			sink(event)
		}
	}, closeTrace, nil
}

func runConformance(stdout io.Writer) int {
	result := conformance.RunSuite(conformance.SyntheticSuite())
	for _, test := range result.Cases {
		status := "PASS"
		if !test.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(stdout, "%s %s steps=%d bdos=%d stop=%s", status, test.Name, test.Steps, test.BDOSCalls, test.StopReason)
		if test.Error != "" {
			fmt.Fprintf(stdout, " error=%s", test.Error)
		}
		fmt.Fprintln(stdout)
	}
	fmt.Fprintf(stdout, "conformance passed=%d failed=%d\n", result.Passed, result.Failed)
	if result.Failed > 0 {
		return 1
	}
	return 0
}

func printDump(w io.Writer, name string, result cpm.RunResult, c *cpu.CPU8080) {
	fmt.Fprintf(w, "\nprogram=%s stop=%s steps=%d bdos=%d\n", name, result.Reason, result.Steps, result.BDOSCalls)
	fmt.Fprintf(w, "A=0x%02X B=0x%02X C=0x%02X D=0x%02X E=0x%02X H=0x%02X L=0x%02X\n", c.A, c.B, c.C, c.D, c.E, c.H, c.L)
	fmt.Fprintf(w, "PC=0x%04X SP=0x%04X Halted=%v Stopped=%v\n", c.PC, c.SP, c.Halted, c.Stopped)
	fmt.Fprintf(w, "Flags C=%v Z=%v S=%v P=%v AC=%v IE=%v\n", c.Carry, c.Zero, c.Sign, c.Parity, c.AuxiliaryCarry, c.InterruptsEnabled)
}
