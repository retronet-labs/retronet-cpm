package cpm

import (
	"errors"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
)

func TestDefaultALUIsNativeAndResetPreservesIt(t *testing.T) {
	m, err := NewMachine(Config{})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	if m.ALUBackend() != cpu.Native {
		t.Fatalf("default ALU is not native")
	}
	m.CPU.SetALU(cpu.Gate)
	m.ResetProgram()
	if m.ALUBackend() != cpu.Native {
		t.Fatalf("machine ALU changed")
	}
}

func TestExplicitGateALU(t *testing.T) {
	m, err := NewMachine(Config{ALU: cpu.Gate})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	if m.ALUBackend() != cpu.Gate {
		t.Fatalf("ALU is not gate")
	}
}

func TestCOMPrintStringAndTerminate(t *testing.T) {
	console := bdos.NewMemoryConsole(nil)
	m, err := NewMachine(Config{Console: console})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	program := []byte{
		cpu.LXI(cpu.PairDE), 0x0D, 0x01,
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		'H', 'I', '$',
	}
	if err := m.LoadCOM("HELLO.COM", program); err != nil {
		t.Fatalf("LoadCOM: %v", err)
	}
	result, err := m.Run(100)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Reason != RunStoppedBDOSTerminate || result.BDOSCalls != 2 || console.Output() != "HI" {
		t.Fatalf("result=%+v output=%q", result, console.Output())
	}
}

func TestWarmBootAtZeroStopsProgram(t *testing.T) {
	m, err := NewMachine(Config{})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	if err := m.LoadCOM("BOOT.COM", []byte{cpu.JMP(), 0x00, 0x00}); err != nil {
		t.Fatalf("LoadCOM: %v", err)
	}
	result, err := m.Run(10)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Reason != RunStoppedWarmBoot {
		t.Fatalf("reason=%s", result.Reason)
	}
}

func TestProgramTooLarge(t *testing.T) {
	m, err := NewMachine(Config{})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	err = m.LoadCOM("BIG.COM", make([]byte, int(BDOSTrapAddress-TransientBase)+1))
	if !errors.Is(err, ErrProgramTooLarge) {
		t.Fatalf("err=%v", err)
	}
}
