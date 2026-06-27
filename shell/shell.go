package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
)

var ErrExit = errors.New("uscita dalla shell")

type Config struct {
	Drive     disk.Drive
	Input     io.Reader
	Output    io.Writer
	ALU       cpu.ALUBackend
	StepLimit uint64
	Trace     cpm.TraceSink
}

type Shell struct {
	drive     disk.Drive
	reader    *bufio.Reader
	output    io.Writer
	console   *bdos.TerminalConsole
	alu       cpu.ALUBackend
	stepLimit uint64
	trace     cpm.TraceSink
}

func New(config Config) (*Shell, error) {
	if config.Drive == nil {
		return nil, errors.New("drive A: non inizializzato")
	}
	input := config.Input
	if input == nil {
		input = strings.NewReader("")
	}
	output := config.Output
	if output == nil {
		output = io.Discard
	}
	alu := config.ALU
	if alu == nil {
		alu = cpu.Native
	}
	stepLimit := config.StepLimit
	if stepLimit == 0 {
		stepLimit = cpm.DefaultStepLimit
	}
	reader := bufio.NewReader(input)
	return &Shell{
		drive:     config.Drive,
		reader:    reader,
		output:    output,
		console:   bdos.NewTerminalConsole(nil, reader, output),
		alu:       alu,
		stepLimit: stepLimit,
		trace:     config.Trace,
	}, nil
}

func (s *Shell) Run() error {
	for {
		fmt.Fprint(s.output, "A>")
		line, err := s.reader.ReadString('\n')
		if err != nil && !(errors.Is(err, io.EOF) && strings.TrimSpace(line) != "") {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if err := s.Execute(line); err != nil {
			if errors.Is(err, ErrExit) {
				return nil
			}
			fmt.Fprintf(s.output, "? %v\n", err)
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
	}
}

func (s *Shell) Execute(line string) error {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return nil
	}
	command := strings.ToUpper(fields[0])
	switch command {
	case "DIR":
		return s.dir()
	case "TYPE":
		if len(fields) != 2 {
			return errors.New("usa TYPE <file>")
		}
		return s.typ(fields[1])
	case "RUN":
		if len(fields) < 2 {
			return errors.New("usa RUN <programma[.COM]> [argomenti]")
		}
		return s.runProgram(fields[1], strings.Join(fields[2:], " "))
	case "HELP":
		fmt.Fprintln(s.output, "DIR  TYPE <file>  RUN <programma[.COM]> [argomenti]  HELP  EXIT")
	case "EXIT":
		return ErrExit
	default:
		return fmt.Errorf("comando sconosciuto: %s", fields[0])
	}
	return nil
}

func (s *Shell) dir() error {
	entries, err := s.drive.List()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		fmt.Fprintf(s.output, "%-12s %6d\n", entry.Name, entry.Size)
	}
	return nil
}

func (s *Shell) typ(name string) error {
	data, err := s.drive.ReadFile(name)
	if err != nil {
		return err
	}
	_, err = s.output.Write(data)
	if err != nil {
		return err
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		fmt.Fprintln(s.output)
	}
	return nil
}

func (s *Shell) runProgram(name string, commandTail string) error {
	name = defaultCOMName(name)
	data, err := s.drive.ReadFile(name)
	if err != nil {
		return err
	}
	m, err := cpm.NewMachine(cpm.Config{
		ALU:       s.alu,
		StepLimit: s.stepLimit,
		Console:   s.console,
		Disk:      s.drive,
		Trace:     s.trace,
	})
	if err != nil {
		return err
	}
	if err := m.LoadCOMWithCommand(name, data, commandTail); err != nil {
		return err
	}
	result, err := m.Run(s.stepLimit)
	if err != nil {
		return err
	}
	fmt.Fprintf(s.output, "\n[%s steps=%d bdos=%d]\n", result.Reason, result.Steps, result.BDOSCalls)
	return nil
}

func defaultCOMName(name string) string {
	if strings.Contains(name, ".") {
		return name
	}
	return name + ".COM"
}
