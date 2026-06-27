// Package conformance contiene programmi sintetici CP/M-like senza ROM storiche.
package conformance

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
	rt "github.com/retronet-labs/retronet-terminal"
)

type Case struct {
	Name       string
	Program    []byte
	Input      []byte
	Command    string
	Disk       disk.Drive
	Terminal   bool
	WantOutput string
	WantReason cpm.RunReason
	WantError  error
	Check      func(*cpm.Machine, string, cpm.RunResult, error) error
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
	mutableDrive := newMemoryDrive()
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
		{
			Name:       "terminal-console-output",
			Program:    consoleOutputProgram('?'),
			Terminal:   true,
			WantOutput: "?",
			WantReason: cpm.RunStoppedBDOSTerminate,
		},
		{
			Name:       "command-tail-default-fcbs",
			Program:    printStringProgram("OK"),
			Command:    "INPUT.TXT OUT.BIN",
			WantOutput: "OK",
			WantReason: cpm.RunStoppedBDOSTerminate,
			Check:      checkCommandTailAndFCBs,
		},
		{
			Name:       "bdos-write-readonly-fails",
			Program:    makeFileOnlyProgram("OUT.TXT"),
			Disk:       readOnlyDrive{},
			WantReason: cpm.RunStoppedBDOSTerminate,
			Check:      checkReturnA(0xFF),
		},
		{
			Name:       "bdos-write-mutable-drive",
			Program:    writeFileProgram("OUT.TXT", []byte("OK$")),
			Disk:       mutableDrive,
			WantReason: cpm.RunStoppedBDOSTerminate,
			Check:      checkWrittenFile(mutableDrive, "OUT.TXT", []byte("OK$")),
		},
	}
}

func RunCase(test Case) CaseResult {
	result := CaseResult{Name: test.Name}
	console, output := newCaseConsole(test)
	m, err := cpm.NewMachine(cpm.Config{Console: console, Disk: test.Disk})
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if err := m.LoadCOMWithCommand(test.Name+".COM", test.Program, test.Command); err != nil {
		result.Error = err.Error()
		return result
	}
	runResult, err := m.Run(cpm.DefaultStepLimit)
	result.Steps = runResult.Steps
	result.BDOSCalls = runResult.BDOSCalls
	result.StopReason = runResult.Reason
	result.Output = output()
	if test.Check != nil {
		if checkErr := test.Check(m, result.Output, runResult, err); checkErr != nil {
			result.Error = checkErr.Error()
			return result
		}
		result.Passed = true
		return result
	}
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

func newCaseConsole(test Case) (bdos.Console, func() string) {
	if test.Terminal {
		term := rt.New(rt.Config{ANSI: true})
		console := bdos.NewTerminalConsole(term, bytes.NewReader(test.Input), nil)
		return console, term.OutputString
	}
	console := bdos.NewMemoryConsole(test.Input)
	return console, console.Output
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

func makeFileOnlyProgram(name string) []byte {
	program := []byte{
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionMakeFile,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
	fcbAddr := cpm.TransientBase + uint16(len(program))
	writeAddr(program, 1, fcbAddr)
	return append(program, fcb(name)...)
}

func writeFileProgram(name string, payload []byte) []byte {
	program := []byte{
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionMakeFile,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionSetDMA,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionWriteSequential,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.LXI(cpu.PairDE), 0x00, 0x00,
		cpu.MVI(cpu.RegC), bdos.FunctionCloseFile,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
	fcbData := fcb(name)
	record := make([]byte, bdos.RecordSize)
	copy(record, payload)
	fcbAddr := cpm.TransientBase + uint16(len(program))
	recordAddr := fcbAddr + uint16(len(fcbData))
	writeAddr(program, 1, fcbAddr)
	writeAddr(program, 9, recordAddr)
	writeAddr(program, 17, fcbAddr)
	writeAddr(program, 25, fcbAddr)
	program = append(program, fcbData...)
	program = append(program, record...)
	return program
}

func fcb(name string) []byte {
	normalized, _ := disk.NormalizeName(name)
	base, ext, _ := strings.Cut(normalized, ".")
	data := make([]byte, 33)
	for i := 1; i <= 11; i++ {
		data[i] = ' '
	}
	copy(data[1:9], []byte(base))
	copy(data[9:12], []byte(ext))
	return data
}

func writeAddr(program []byte, offset int, addr uint16) {
	program[offset] = byte(addr)
	program[offset+1] = byte(addr >> 8)
}

func checkCommandTailAndFCBs(m *cpm.Machine, output string, runResult cpm.RunResult, err error) error {
	if err != nil {
		return err
	}
	if runResult.Reason != cpm.RunStoppedBDOSTerminate || output != "OK" {
		return fmt.Errorf("reason=%s output=%q", runResult.Reason, output)
	}
	tail := "INPUT.TXT OUT.BIN"
	if got := m.Memory.Read(cpm.CommandTailAddr); got != byte(len(tail)) {
		return fmt.Errorf("tail len=%d want=%d", got, len(tail))
	}
	for i, value := range []byte(tail) {
		if got := m.Memory.Read(cpm.CommandTailAddr + 1 + uint16(i)); got != value {
			return fmt.Errorf("tail[%d]=0x%02X want=0x%02X", i, got, value)
		}
	}
	if err := checkFCB(m, cpm.DefaultFCB1, "INPUT   ", "TXT"); err != nil {
		return err
	}
	return checkFCB(m, cpm.DefaultFCB2, "OUT     ", "BIN")
}

func checkFCB(m *cpm.Machine, addr uint16, name string, ext string) error {
	for i := 0; i < 8; i++ {
		if got := m.Memory.Read(addr + 1 + uint16(i)); got != name[i] {
			return fmt.Errorf("fcb name[%d]=0x%02X want=0x%02X", i, got, name[i])
		}
	}
	for i := 0; i < 3; i++ {
		if got := m.Memory.Read(addr + 9 + uint16(i)); got != ext[i] {
			return fmt.Errorf("fcb ext[%d]=0x%02X want=0x%02X", i, got, ext[i])
		}
	}
	return nil
}

func checkReturnA(want byte) func(*cpm.Machine, string, cpm.RunResult, error) error {
	return func(m *cpm.Machine, _ string, runResult cpm.RunResult, err error) error {
		if err != nil {
			return err
		}
		if runResult.Reason != cpm.RunStoppedBDOSTerminate {
			return fmt.Errorf("reason=%s", runResult.Reason)
		}
		if m.CPU.A != want {
			return fmt.Errorf("A=0x%02X want=0x%02X", m.CPU.A, want)
		}
		return nil
	}
}

func checkWrittenFile(drive *memoryDrive, name string, prefix []byte) func(*cpm.Machine, string, cpm.RunResult, error) error {
	return func(_ *cpm.Machine, _ string, runResult cpm.RunResult, err error) error {
		if err != nil {
			return err
		}
		if runResult.Reason != cpm.RunStoppedBDOSTerminate {
			return fmt.Errorf("reason=%s", runResult.Reason)
		}
		data, ok := drive.files[name]
		if !ok {
			return fmt.Errorf("%s non scritto", name)
		}
		if !bytes.HasPrefix(data, prefix) {
			return fmt.Errorf("%s prefix=%q want=%q", name, data[:len(prefix)], prefix)
		}
		return nil
	}
}

type readOnlyDrive struct{}

func (readOnlyDrive) List() ([]disk.Entry, error) { return nil, nil }

func (readOnlyDrive) ReadFile(string) ([]byte, error) { return nil, os.ErrNotExist }

type memoryDrive struct {
	files map[string][]byte
}

func newMemoryDrive() *memoryDrive {
	return &memoryDrive{files: map[string][]byte{}}
}

func (d *memoryDrive) List() ([]disk.Entry, error) {
	entries := make([]disk.Entry, 0, len(d.files))
	for name, data := range d.files {
		entries = append(entries, disk.Entry{Name: name, Size: int64(len(data))})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries, nil
}

func (d *memoryDrive) ReadFile(name string) ([]byte, error) {
	data, ok := d.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return append([]byte(nil), data...), nil
}

func (d *memoryDrive) WriteFile(name string, data []byte) error {
	normalized, err := disk.NormalizeName(name)
	if err != nil {
		return err
	}
	d.files[normalized] = append([]byte(nil), data...)
	return nil
}

func (d *memoryDrive) DeleteFile(name string) error {
	normalized, err := disk.NormalizeName(name)
	if err != nil {
		return err
	}
	delete(d.files, normalized)
	return nil
}

func (d *memoryDrive) RenameFile(oldName string, newName string) error {
	oldNormalized, err := disk.NormalizeName(oldName)
	if err != nil {
		return err
	}
	newNormalized, err := disk.NormalizeName(newName)
	if err != nil {
		return err
	}
	data, ok := d.files[oldNormalized]
	if !ok {
		return os.ErrNotExist
	}
	delete(d.files, oldNormalized)
	d.files[newNormalized] = data
	return nil
}
