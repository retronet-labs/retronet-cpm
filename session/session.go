// Package session espone un contratto programmatico API-ready sopra la shell
// CP/M-like, senza passare dalla CLI.
package session

import (
	"errors"
	"io"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
	"github.com/retronet-labs/retronet-cpm/shell"
	rt "github.com/retronet-labs/retronet-terminal"
)

var ErrNilSession = errors.New("sessione CP/M non inizializzata")

type Config struct {
	Drive     disk.Drive
	ALU       cpu.ALUBackend
	StepLimit uint64
	Trace     cpm.TraceSink
	Terminal  *rt.Terminal
	Input     io.Reader
	Output    io.Writer
}

type Session struct {
	shell    *shell.Shell
	terminal *rt.Terminal
}

func New(config Config) (*Session, error) {
	terminal := config.Terminal
	if terminal == nil {
		terminal = rt.New(rt.Config{ANSI: true})
	}
	output := config.Output
	if output == nil {
		output = terminal
	}
	sh, err := shell.New(shell.Config{
		Drive:     config.Drive,
		Input:     config.Input,
		Output:    output,
		ALU:       config.ALU,
		StepLimit: config.StepLimit,
		Trace:     config.Trace,
		Terminal:  terminal,
	})
	if err != nil {
		return nil, err
	}
	return &Session{shell: sh, terminal: terminal}, nil
}

func (s *Session) Input(data []byte) error {
	if s == nil || s.terminal == nil {
		return ErrNilSession
	}
	s.terminal.QueueInput(data)
	return nil
}

func (s *Session) RunCommand(line string) error {
	if s == nil || s.shell == nil {
		return ErrNilSession
	}
	return s.shell.Execute(line)
}

func (s *Session) Prompt() error {
	if s == nil || s.terminal == nil {
		return ErrNilSession
	}
	_, err := s.terminal.Write([]byte("A>"))
	return err
}

func (s *Session) DrainOutput() ([]byte, error) {
	if s == nil || s.terminal == nil {
		return nil, ErrNilSession
	}
	return s.terminal.DrainOutput(), nil
}

func (s *Session) Snapshot() (rt.Snapshot, error) {
	if s == nil || s.terminal == nil {
		return rt.Snapshot{}, ErrNilSession
	}
	return s.terminal.Snapshot(), nil
}

func (s *Session) Terminal() *rt.Terminal {
	if s == nil {
		return nil
	}
	return s.terminal
}
