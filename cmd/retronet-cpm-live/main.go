// Comando retronet-cpm-live: shell CP/M-like interattiva su retronet-terminal.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
	"github.com/retronet-labs/retronet-cpm/session"
	"github.com/retronet-labs/retronet-cpm/shell"
	rt "github.com/retronet-labs/retronet-terminal"
	"github.com/retronet-labs/retronet-terminal/live"
)

const footer = "A> live | Invio esegue | Ctrl+L pulisce | Ctrl+Q/Ctrl+C esce"

type runConfig struct {
	diskPath   string
	steps      uint64
	aluName    string
	writeDisk  bool
	width      int
	height     int
	lineMode   bool
	scriptMode bool
	script     string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, true))
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, raw bool) int {
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
	drive, err := openDrive(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "errore drive A: %v\n", err)
		return 2
	}
	handler := &cpmHandler{
		drive:     drive,
		alu:       alu,
		stepLimit: cfg.steps,
	}
	err = live.Run(live.Config{
		Width:      cfg.width,
		Height:     cfg.height,
		Input:      stdin,
		Output:     stdout,
		Raw:        raw,
		LineMode:   cfg.lineMode,
		ScriptMode: cfg.scriptMode,
		Script:     []byte(cfg.script),
		Footer:     footer,
		Handler:    handler,
	})
	if err != nil {
		fmt.Fprintf(stderr, "errore terminale CP/M live: %v\n", err)
		if !cfg.lineMode {
			fmt.Fprintln(stderr, "Riprova da una console interattiva, oppure usa -line per la modalita' a righe.")
		}
		return 1
	}
	return 0
}

func parseFlags(args []string, stderr io.Writer) (runConfig, error) {
	fs := flag.NewFlagSet("retronet-cpm-live", flag.ContinueOnError)
	fs.SetOutput(stderr)
	cfg := runConfig{
		diskPath: ".",
		steps:    cpm.DefaultStepLimit,
		aluName:  defaultALUName(),
		width:    rt.DefaultWidth,
		height:   rt.DefaultHeight,
	}
	fs.StringVar(&cfg.diskPath, "disk", cfg.diskPath, "directory host mappata come drive A:")
	fs.Uint64Var(&cfg.steps, "steps", cfg.steps, "limite massimo di istruzioni 8080 per RUN")
	fs.StringVar(&cfg.aluName, "alu", cfg.aluName, "backend ALU: native o gate")
	fs.BoolVar(&cfg.writeDisk, "write-disk", false, "abilita funzioni BDOS che modificano il drive host")
	fs.IntVar(&cfg.width, "width", cfg.width, "larghezza dello schermo")
	fs.IntVar(&cfg.height, "height", cfg.height, "altezza dello schermo")
	fs.BoolVar(&cfg.lineMode, "line", false, "usa input a righe invece del raw mode")
	fs.StringVar(&cfg.script, "script", "", "comandi da inviare al terminale live e poi terminare")
	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	cfg.scriptMode = cfg.script != ""
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

type cpmHandler struct {
	drive     disk.Drive
	alu       cpu.ALUBackend
	stepLimit uint64
	session   *session.Session
	line      []byte
}

func (h *cpmHandler) Start(term *rt.Terminal) error {
	sess, err := session.New(session.Config{
		Drive:     h.drive,
		ALU:       h.alu,
		StepLimit: h.stepLimit,
		Terminal:  term,
		Output:    term,
	})
	if err != nil {
		return err
	}
	h.session = sess
	if _, err := term.Write([]byte("RetroNet CP/M Live\r\n")); err != nil {
		return err
	}
	return h.session.Prompt()
}

func (h *cpmHandler) HandleByte(term *rt.Terminal, value byte) (bool, error) {
	switch value {
	case 0x03, 0x04, 0x11: // Ctrl+C, Ctrl+D, Ctrl+Q.
		_, err := term.Write([]byte("\r\n"))
		return false, err
	case 0x0C: // Ctrl+L.
		h.line = h.line[:0]
		if _, err := term.Write([]byte("\x1b[2J\x1b[H")); err != nil {
			return true, err
		}
		return true, h.session.Prompt()
	case '\r', '\n':
		return h.executeLine(term)
	case '\b', 0x7F:
		if len(h.line) == 0 {
			return true, nil
		}
		h.line = h.line[:len(h.line)-1]
		_, err := term.Write([]byte{'\b', ' ', '\b'})
		return true, err
	default:
		if value == '\t' || (value >= 0x20 && value <= 0x7E) {
			h.line = append(h.line, value)
			return true, term.WriteByte(value)
		}
	}
	return true, nil
}

func (h *cpmHandler) executeLine(term *rt.Terminal) (bool, error) {
	command := strings.TrimSpace(string(h.line))
	h.line = h.line[:0]
	if _, err := term.Write([]byte("\r\n")); err != nil {
		return true, err
	}
	if command == "" {
		return true, h.session.Prompt()
	}
	if err := h.session.RunCommand(command); err != nil {
		if errors.Is(err, shell.ErrExit) {
			return false, nil
		}
		if _, writeErr := fmt.Fprintf(term, "? %v\r\n", err); writeErr != nil {
			return true, writeErr
		}
	}
	return true, h.session.Prompt()
}
