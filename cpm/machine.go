// Package cpm modella una macchina CP/M-like sopra l'8080 RetroNet.
package cpm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/disk"
)

const (
	TransientBase    uint16 = 0x0100
	BDOSVector       uint16 = 0x0005
	BDOSTrapAddress  uint16 = 0xF000
	DefaultStack     uint16 = 0xEFFE
	DefaultStepLimit uint64 = 1_000_000
	CommandTailAddr  uint16 = 0x0080
	DefaultFCB1      uint16 = 0x005C
	DefaultFCB2      uint16 = 0x006C
	MaxCommandTail   int    = 126
)

var ErrProgramTooLarge = fmt.Errorf("programma COM troppo grande")
var ErrCommandTailTooLong = errors.New("command tail CP/M troppo lungo")

// Config contiene le dipendenze della macchina CP/M-like.
type Config struct {
	ALU       cpu.ALUBackend
	StepLimit uint64
	Console   bdos.Console
	Disk      disk.Drive
	Trace     TraceSink
}

type RunReason string

const (
	RunStoppedHalted        RunReason = "halted"
	RunStoppedWarmBoot      RunReason = "warm-boot"
	RunStoppedBDOSTerminate RunReason = "bdos-terminate"
	RunStoppedLimit         RunReason = "step-limit"
)

type RunResult struct {
	Reason    RunReason `json:"reason"`
	Steps     uint64    `json:"steps"`
	BDOSCalls uint64    `json:"bdos_calls"`
	PC        uint16    `json:"pc"`
	SP        uint16    `json:"sp"`
}

type TraceKind string

const (
	TraceInstruction TraceKind = "instruction"
	TraceBDOS        TraceKind = "bdos"
)

type TraceEvent struct {
	Sequence     uint64      `json:"sequence"`
	Kind         TraceKind   `json:"kind"`
	PC           uint16      `json:"pc"`
	Opcode       byte        `json:"opcode,omitempty"`
	Disassembly  string      `json:"disassembly,omitempty"`
	BDOSFunction byte        `json:"bdos_function,omitempty"`
	Before       cpu.CPU8080 `json:"before"`
	After        cpu.CPU8080 `json:"after"`
}

type TraceSink func(TraceEvent)

// Machine e' una macchina CP/M-like isolata.
type Machine struct {
	CPU    *cpu.CPU8080
	Memory *cpu.FlatMemory
	IO     *cpu.Ports

	config        Config
	alu           cpu.ALUBackend
	bdos          *bdos.Handler
	traceSequence uint64
	loadedName    string
}

func NewMachine(config Config) (*Machine, error) {
	alu := config.ALU
	if alu == nil {
		alu = cpu.Native
	}
	stepLimit := config.StepLimit
	if stepLimit == 0 {
		stepLimit = DefaultStepLimit
	}
	console := config.Console
	if console == nil {
		console = bdos.NewMemoryConsole(nil)
	}
	config.ALU = alu
	config.StepLimit = stepLimit
	config.Console = console

	m := &Machine{
		CPU:    cpu.NewCPU8080WithALU(alu),
		Memory: cpu.NewFlatMemory(),
		IO:     cpu.NewPorts(),
		config: config,
		alu:    alu,
		bdos:   bdos.NewHandler(console),
	}
	m.bdos.Disk = config.Disk
	m.ResetProgram()
	return m, nil
}

func (m *Machine) ALUBackend() cpu.ALUBackend { return m.alu }

func (m *Machine) ResetProgram() {
	m.Memory.Data = [cpu.AddressSpaceSize]byte{}
	m.CPU.Reset()
	m.CPU.PC = TransientBase
	m.CPU.SP = DefaultStack
	if m.bdos != nil {
		m.bdos.Reset()
		m.bdos.Disk = m.config.Disk
	}
	m.installPageZero()
	m.traceSequence = 0
	m.loadedName = ""
}

func (m *Machine) installPageZero() {
	m.Memory.Write(0x0000, cpu.JMP())
	m.Memory.Write(0x0001, 0x00)
	m.Memory.Write(0x0002, 0x00)
	m.Memory.Write(BDOSVector, cpu.JMP())
	m.Memory.Write(BDOSVector+1, byte(BDOSTrapAddress&0x00FF))
	m.Memory.Write(BDOSVector+2, byte(BDOSTrapAddress>>8))
	m.Memory.Write(BDOSTrapAddress, cpu.RET())
}

func (m *Machine) LoadCOM(name string, data []byte) error {
	if len(data) > int(BDOSTrapAddress-TransientBase) {
		return fmt.Errorf("%w: %d byte, max %d", ErrProgramTooLarge, len(data), BDOSTrapAddress-TransientBase)
	}
	m.ResetProgram()
	for i, value := range data {
		m.Memory.Write(TransientBase+uint16(i), value)
	}
	m.CPU.PC = TransientBase
	m.CPU.SP = DefaultStack
	m.loadedName = name
	return nil
}

func (m *Machine) LoadCOMWithCommand(name string, data []byte, commandTail string) error {
	if err := m.LoadCOM(name, data); err != nil {
		return err
	}
	if err := m.SetCommandTail(commandTail); err != nil {
		return err
	}
	return m.SetDefaultFCBs(commandTail)
}

func (m *Machine) SetCommandTail(commandTail string) error {
	if len(commandTail) > MaxCommandTail {
		return fmt.Errorf("%w: %d byte, max %d", ErrCommandTailTooLong, len(commandTail), MaxCommandTail)
	}
	m.Memory.Write(CommandTailAddr, byte(len(commandTail)))
	for i := 0; i < MaxCommandTail; i++ {
		value := byte(0)
		if i < len(commandTail) {
			value = commandTail[i]
		}
		m.Memory.Write(CommandTailAddr+1+uint16(i), value)
	}
	m.Memory.Write(CommandTailAddr+1+uint16(len(commandTail)), '\r')
	return nil
}

func (m *Machine) SetDefaultFCBs(commandTail string) error {
	clearFCB(m.Memory, DefaultFCB1)
	clearFCB(m.Memory, DefaultFCB2)
	fields := strings.Fields(commandTail)
	if len(fields) > 0 {
		if err := writeFCB(m.Memory, DefaultFCB1, fields[0]); err != nil {
			return err
		}
	}
	if len(fields) > 1 {
		if err := writeFCB(m.Memory, DefaultFCB2, fields[1]); err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) Run(limit uint64) (RunResult, error) {
	if limit == 0 {
		limit = m.config.StepLimit
	}
	var result RunResult
	for result.Steps < limit {
		if m.CPU.Halted || m.CPU.Stopped {
			return m.finish(result, RunStoppedHalted), nil
		}
		switch m.CPU.PC {
		case 0x0000:
			return m.finish(result, RunStoppedWarmBoot), nil
		case BDOSVector, BDOSTrapAddress:
			if err := m.callBDOS(&result); err != nil {
				return m.finish(result, ""), err
			}
			if m.CPU.Halted || m.CPU.Stopped {
				return m.finish(result, RunStoppedBDOSTerminate), nil
			}
			continue
		}
		if err := m.stepInstruction(&result); err != nil {
			return m.finish(result, ""), err
		}
	}
	return m.finish(result, RunStoppedLimit), nil
}

func (m *Machine) stepInstruction(result *RunResult) error {
	before := *m.CPU
	disassembly, err := cpu.Disassemble(m.Memory, m.CPU.PC)
	if err != nil {
		return err
	}
	if err := m.CPU.Step(m.Memory, m.IO); err != nil {
		return err
	}
	result.Steps++
	m.emit(TraceEvent{
		Kind:        TraceInstruction,
		PC:          before.PC,
		Opcode:      disassembly.Opcode.Code,
		Disassembly: disassembly.String(),
		Before:      before,
		After:       *m.CPU,
	})
	return nil
}

func (m *Machine) callBDOS(result *RunResult) error {
	before := *m.CPU
	callResult, err := m.bdos.Call(m.CPU, m.Memory)
	if err != nil {
		return err
	}
	result.BDOSCalls++
	if callResult.Terminate {
		m.CPU.Halted = true
		m.CPU.Stopped = true
	} else {
		m.returnFromCall()
	}
	m.emit(TraceEvent{
		Kind:         TraceBDOS,
		PC:           before.PC,
		BDOSFunction: callResult.Function,
		Before:       before,
		After:        *m.CPU,
	})
	return nil
}

func (m *Machine) returnFromCall() {
	low := uint16(m.Memory.Read(m.CPU.SP))
	high := uint16(m.Memory.Read(m.CPU.SP + 1))
	m.CPU.SP += 2
	m.CPU.PC = high<<8 | low
}

func (m *Machine) emit(event TraceEvent) {
	if m.config.Trace == nil {
		return
	}
	event.Sequence = m.traceSequence
	m.traceSequence++
	m.config.Trace(event)
}

func (m *Machine) finish(result RunResult, reason RunReason) RunResult {
	result.Reason = reason
	result.PC = m.CPU.PC
	result.SP = m.CPU.SP
	return result
}

func clearFCB(mem *cpu.FlatMemory, addr uint16) {
	mem.Write(addr, 0)
	for i := 1; i <= 11; i++ {
		mem.Write(addr+uint16(i), ' ')
	}
	for i := 12; i < 16; i++ {
		mem.Write(addr+uint16(i), 0)
	}
}

func writeFCB(mem *cpu.FlatMemory, addr uint16, name string) error {
	normalized, err := disk.NormalizeName(name)
	if err != nil {
		return err
	}
	base, ext, _ := strings.Cut(normalized, ".")
	mem.Write(addr, 0)
	for i := 0; i < 8; i++ {
		value := byte(' ')
		if i < len(base) {
			value = base[i]
		}
		mem.Write(addr+1+uint16(i), value)
	}
	for i := 0; i < 3; i++ {
		value := byte(' ')
		if i < len(ext) {
			value = ext[i]
		}
		mem.Write(addr+9+uint16(i), value)
	}
	return nil
}
