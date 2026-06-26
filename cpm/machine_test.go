package cpm

import (
	"errors"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/disk"
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

func TestCOMReadsFCBFileThroughBDOS(t *testing.T) {
	console := bdos.NewMemoryConsole(nil)
	m, err := NewMachine(Config{
		Console: console,
		Disk:    fakeDisk{files: map[string][]byte{"MSG.TXT": []byte("OK$")}},
	})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	if err := m.LoadCOM("READ.COM", fcbReadProgram()); err != nil {
		t.Fatalf("LoadCOM: %v", err)
	}
	result, err := m.Run(100)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Reason != RunStoppedBDOSTerminate || console.Output() != "OK" {
		t.Fatalf("result=%+v output=%q", result, console.Output())
	}
}

func fcbReadProgram() []byte {
	const dmaAddr = 0x0200
	program := []byte{
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionOpenFile,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), byte(dmaAddr & 0x00FF), byte(dmaAddr >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionSetDMA,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionReadSequential,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), byte(dmaAddr & 0x00FF), byte(dmaAddr >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(BDOSVector), byte(BDOSVector >> 8),
	}
	fcbAddr := TransientBase + uint16(len(program))
	program[1] = byte(fcbAddr)
	program[2] = byte(fcbAddr >> 8)
	program[17] = byte(fcbAddr)
	program[18] = byte(fcbAddr >> 8)
	fcb := make([]byte, 33)
	copy(fcb[1:9], []byte("MSG     "))
	copy(fcb[9:12], []byte("TXT"))
	return append(program, fcb...)
}

type fakeDisk struct {
	files map[string][]byte
}

func (d fakeDisk) List() ([]disk.Entry, error) {
	entries := make([]disk.Entry, 0, len(d.files))
	for name, data := range d.files {
		entries = append(entries, disk.Entry{Name: name, Size: int64(len(data))})
	}
	return entries, nil
}

func (d fakeDisk) ReadFile(name string) ([]byte, error) {
	data, ok := d.files[name]
	if !ok {
		return nil, errors.New("missing")
	}
	return append([]byte(nil), data...), nil
}
