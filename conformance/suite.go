// Package conformance contiene programmi sintetici CP/M-like senza ROM storiche.
package conformance

import (
	"errors"
	"fmt"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
)

type Case struct {
	Name       string
	Program    []byte
	Input      []byte
	WantOutput string
	WantReason cpm.RunReason
	WantError  error
}

type CaseResult struct {
	Name       string        `json:"name"`
	Passed     bool          `json:"passed"`
	Steps      uint64        `json:"steps"`
	BDOSCalls  uint64        `json:"bdos_calls"`
	StopReason cpm.RunReason `json:"stop_reason"`
	Error      string        `json:"error,omitempty"`
	Output     string        `json:"output,omitempty"`
}

type SuiteResult struct {
	Passed int          `json:"passed"`
	Failed int          `json:"failed"`
	Cases  []CaseResult `json:"cases"`
}

func SyntheticSuite() []Case {
	return []Case{
		{
			Name:       "bdos-print-string",
			Program:    printStringProgram("HELLO"),
			WantOutput: "HELLO",
			WantReason: cpm.RunStoppedBDOSTerminate,
		},
		{
			Name:       "bdos-console-output",
			Program:    consoleOutputProgram('!'),
			WantOutput: "!",
			WantReason: cpm.RunStoppedBDOSTerminate,
		},
		{
			Name:       "bdos-direct-console-input",
			Program:    directInputProgram(),
			Input:      []byte("Z"),
			WantOutput: "Z",
			WantReason: cpm.RunStoppedBDOSTerminate,
		},
		{
			Name:       "warm-boot",
			Program:    []byte{cpu.JMP(), 0x00, 0x00},
			WantReason: cpm.RunStoppedWarmBoot,
		},
		{
			Name:      "unsupported-bdos",
			Program:   unsupportedProgram(),
			WantError: bdos.ErrUnsupportedFunction,
		},
	}
}

func RunCase(test Case) CaseResult {
	result := CaseResult{Name: test.Name}
	console := bdos.NewMemoryConsole(test.Input)
	m, err := cpm.NewMachine(cpm.Config{Console: console})
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if err := m.LoadCOM(test.Name+".COM", test.Program); err != nil {
		result.Error = err.Error()
		return result
	}
	runResult, err := m.Run(cpm.DefaultStepLimit)
	result.Steps = runResult.Steps
	result.BDOSCalls = runResult.BDOSCalls
	result.StopReason = runResult.Reason
	result.Output = console.Output()
	if test.WantError != nil {
		if errors.Is(err, test.WantError) {
			result.Passed = true
			return result
		}
		result.Error = fmt.Sprintf("errore=%v want=%v", err, test.WantError)
		return result
	}
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if test.WantReason != "" && runResult.Reason != test.WantReason {
		result.Error = fmt.Sprintf("stop=%s want=%s", runResult.Reason, test.WantReason)
		return result
	}
	if result.Output != test.WantOutput {
		result.Error = fmt.Sprintf("output=%q want=%q", result.Output, test.WantOutput)
		return result
	}
	result.Passed = true
	return result
}

func RunSuite(cases []Case) SuiteResult {
	result := SuiteResult{Cases: make([]CaseResult, 0, len(cases))}
	for _, test := range cases {
		caseResult := RunCase(test)
		result.Cases = append(result.Cases, caseResult)
		if caseResult.Passed {
			result.Passed++
		} else {
			result.Failed++
		}
	}
	return result
}

func printStringProgram(text string) []byte {
	program := []byte{
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
	msgAddr := cpm.TransientBase + uint16(len(program))
	program[1] = byte(msgAddr)
	program[2] = byte(msgAddr >> 8)
	program = append(program, []byte(text)...)
	program = append(program, '$')
	return program
}

func consoleOutputProgram(value byte) []byte {
	return []byte{
		cpu.MVI(cpu.RegE), value,
		cpu.MVI(cpu.RegC), bdos.FunctionConsoleOutput,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
}

func directInputProgram() []byte {
	return []byte{
		cpu.MVI(cpu.RegE), 0xFF,
		cpu.MVI(cpu.RegC), bdos.FunctionDirectConsoleIO,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MOV(cpu.RegE, cpu.RegA),
		cpu.MVI(cpu.RegC), bdos.FunctionConsoleOutput,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
}

func unsupportedProgram() []byte {
	return []byte{
		cpu.MVI(cpu.RegC), 99,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
}
