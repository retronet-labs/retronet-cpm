package bdos

import (
	"bytes"
	"strings"
	"testing"

	rt "github.com/retronet-labs/retronet-terminal"
)

func TestTerminalConsoleUsesQueuedInputBeforeReader(t *testing.T) {
	term := rt.New(rt.Config{ANSI: true})
	console := NewTerminalConsole(term, strings.NewReader("B"), nil)
	console.QueueInput([]byte("A"))

	a, err := console.ReadByte()
	if err != nil || a != 'A' {
		t.Fatalf("first input=0x%02X err=%v", a, err)
	}
	b, err := console.ReadByte()
	if err != nil || b != 'B' {
		t.Fatalf("second input=0x%02X err=%v", b, err)
	}
}

func TestTerminalConsoleBuffersAndMirrorsOutput(t *testing.T) {
	var mirror bytes.Buffer
	term := rt.New(rt.Config{Width: 8, Height: 2, ANSI: true})
	console := NewTerminalConsole(term, nil, &mirror)

	if err := console.WriteByte('O'); err != nil {
		t.Fatal(err)
	}
	if err := console.WriteByte('K'); err != nil {
		t.Fatal(err)
	}
	if mirror.String() != "OK" || term.OutputString() != "OK" {
		t.Fatalf("mirror=%q terminal=%q", mirror.String(), term.OutputString())
	}
	if got := term.ScreenString(); got != "OK\n" {
		t.Fatalf("screen=%q", got)
	}
}
