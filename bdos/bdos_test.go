package bdos

import (
	"errors"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
)

func TestPrintStringAndTerminate(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	c.C = FunctionPrintString
	c.SetDE(0x0200)
	for i, value := range []byte("CIAO$") {
		mem.Write(0x0200+uint16(i), value)
	}
	console := NewMemoryConsole(nil)
	handler := NewHandler(console)

	result, err := handler.Call(c, mem)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if result.Terminate || console.Output() != "CIAO" {
		t.Fatalf("result=%+v output=%q", result, console.Output())
	}

	c.C = FunctionTerminate
	result, err = handler.Call(c, mem)
	if err != nil {
		t.Fatalf("terminate: %v", err)
	}
	if !result.Terminate {
		t.Fatalf("terminate=false")
	}
}

func TestConsoleFunctions(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	console := NewMemoryConsole([]byte("Aline\r"))
	handler := NewHandler(console)

	c.C = FunctionConsoleStatus
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("status: %v", err)
	}
	if c.A != 0xFF {
		t.Fatalf("status A=0x%02X", c.A)
	}

	c.C = FunctionConsoleInput
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("input: %v", err)
	}
	if c.A != 'A' || console.Output() != "A" {
		t.Fatalf("A=0x%02X output=%q", c.A, console.Output())
	}

	c.C = FunctionConsoleOutput
	c.E = '!'
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("output: %v", err)
	}
	if console.Output() != "A!" {
		t.Fatalf("output=%q", console.Output())
	}

	c.C = FunctionReadConsoleLine
	c.SetDE(0x0300)
	mem.Write(0x0300, 8)
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("line: %v", err)
	}
	if mem.Read(0x0301) != 4 || string([]byte{mem.Read(0x0302), mem.Read(0x0303), mem.Read(0x0304), mem.Read(0x0305)}) != "line" {
		t.Fatalf("count=%d data=%q", mem.Read(0x0301), []byte{mem.Read(0x0302), mem.Read(0x0303), mem.Read(0x0304), mem.Read(0x0305)})
	}
}

func TestUnsupportedFunction(t *testing.T) {
	c := cpu.NewCPU8080()
	c.C = 99
	_, err := NewHandler(NewMemoryConsole(nil)).Call(c, cpu.NewFlatMemory())
	if !errors.Is(err, ErrUnsupportedFunction) {
		t.Fatalf("err=%v", err)
	}
}
