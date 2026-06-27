package bdos

import (
	"bufio"
	"io"
	"strings"

	rt "github.com/retronet-labs/retronet-terminal"
)

// TerminalConsole adatta retronet-terminal all'interfaccia Console del BDOS.
// Mantiene il buffer/schermo del terminale e, opzionalmente, rispecchia l'output
// su uno writer host per preservare il comportamento della CLI.
type TerminalConsole struct {
	terminal *rt.Terminal
	reader   *bufio.Reader
	mirror   io.Writer
}

func NewTerminalConsole(term *rt.Terminal, input io.Reader, mirror io.Writer) *TerminalConsole {
	if term == nil {
		term = rt.New(rt.Config{ANSI: true})
	}
	if input == nil {
		input = strings.NewReader("")
	}
	if mirror == nil {
		mirror = io.Discard
	}
	return &TerminalConsole{
		terminal: term,
		reader:   bufio.NewReader(input),
		mirror:   mirror,
	}
}

func (c *TerminalConsole) Terminal() *rt.Terminal {
	return c.terminal
}

func (c *TerminalConsole) QueueInput(data []byte) {
	c.terminal.QueueInput(data)
}

func (c *TerminalConsole) ReadByte() (byte, error) {
	if c.terminal.Status() {
		return c.terminal.ReadByte()
	}
	return c.reader.ReadByte()
}

func (c *TerminalConsole) WriteByte(value byte) error {
	if err := c.terminal.WriteByte(value); err != nil {
		return err
	}
	_, err := c.mirror.Write([]byte{value})
	return err
}

func (c *TerminalConsole) Status() bool {
	return c.terminal.Status() || c.reader.Buffered() > 0
}
