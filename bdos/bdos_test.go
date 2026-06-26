package bdos

import (
	"errors"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/disk"
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

func TestFCBReadOnlySequential(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	fcbAddr := uint16(0x005C)
	dmaAddr := uint16(0x0200)
	writeFCB(mem, fcbAddr, "README", "TXT")
	handler := NewHandler(NewMemoryConsole(nil))
	handler.Disk = fakeDrive{files: map[string][]byte{"README.TXT": []byte("hello")}}

	c.SetDE(dmaAddr)
	c.C = FunctionSetDMA
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("set dma: %v", err)
	}
	if handler.DMA != dmaAddr {
		t.Fatalf("DMA=0x%04X", handler.DMA)
	}

	c.SetDE(fcbAddr)
	c.C = FunctionOpenFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("open: %v", err)
	}
	if c.A != 0 || mem.Read(fcbAddr+32) != 0 {
		t.Fatalf("open A=0x%02X CR=%d", c.A, mem.Read(fcbAddr+32))
	}

	c.C = FunctionReadSequential
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("read: %v", err)
	}
	got := []byte{mem.Read(dmaAddr), mem.Read(dmaAddr + 1), mem.Read(dmaAddr + 2), mem.Read(dmaAddr + 3), mem.Read(dmaAddr + 4), mem.Read(dmaAddr + 5)}
	if string(got[:5]) != "hello" || got[5] != 0x1A || c.A != 0 || mem.Read(fcbAddr+32) != 1 {
		t.Fatalf("read got=%v A=0x%02X CR=%d", got, c.A, mem.Read(fcbAddr+32))
	}

	c.C = FunctionReadSequential
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("eof read: %v", err)
	}
	if c.A != 1 {
		t.Fatalf("EOF A=0x%02X", c.A)
	}

	c.C = FunctionCloseFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("close: %v", err)
	}
	if c.A != 0 {
		t.Fatalf("close A=0x%02X", c.A)
	}
}

func TestFCBOpenWithoutDisk(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	writeFCB(mem, 0x005C, "MISSING", "TXT")
	c.SetDE(0x005C)
	c.C = FunctionOpenFile
	_, err := NewHandler(NewMemoryConsole(nil)).Call(c, mem)
	if !errors.Is(err, ErrNilDisk) {
		t.Fatalf("err=%v", err)
	}
}

func TestFCBWriteDeleteRenameRequiresMutableDrive(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	fcbAddr := uint16(0x005C)
	dmaAddr := uint16(0x0200)
	writeFCB(mem, fcbAddr, "OUT", "TXT")
	for i, value := range []byte("created$") {
		mem.Write(dmaAddr+uint16(i), value)
	}
	handler := NewHandler(NewMemoryConsole(nil))
	handler.Disk = fakeDrive{files: map[string][]byte{}}

	c.SetDE(dmaAddr)
	c.C = FunctionSetDMA
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("set dma: %v", err)
	}
	c.SetDE(fcbAddr)
	c.C = FunctionMakeFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("make: %v", err)
	}
	if c.A != 0 {
		t.Fatalf("make A=0x%02X", c.A)
	}
	c.C = FunctionWriteSequential
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("write: %v", err)
	}
	if c.A != 0 {
		t.Fatalf("write A=0x%02X", c.A)
	}
	c.C = FunctionCloseFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("close: %v", err)
	}
	data := handler.Disk.(fakeDrive).files["OUT.TXT"]
	if string(data[:8]) != "created$" {
		t.Fatalf("written=%q", data[:8])
	}

	writeRenameFCB(mem, fcbAddr, "NEW", "TXT")
	c.C = FunctionRenameFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("rename: %v", err)
	}
	if _, ok := handler.Disk.(fakeDrive).files["NEW.TXT"]; !ok || c.A != 0 {
		t.Fatalf("rename A=0x%02X files=%v", c.A, handler.Disk.(fakeDrive).files)
	}
	writeFCB(mem, fcbAddr, "NEW", "TXT")
	c.C = FunctionDeleteFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := handler.Disk.(fakeDrive).files["NEW.TXT"]; ok || c.A != 0 {
		t.Fatalf("delete A=0x%02X files=%v", c.A, handler.Disk.(fakeDrive).files)
	}
}

func TestFCBWriteFailsOnReadOnlyDrive(t *testing.T) {
	mem := cpu.NewFlatMemory()
	c := cpu.NewCPU8080()
	writeFCB(mem, 0x005C, "OUT", "TXT")
	handler := NewHandler(NewMemoryConsole(nil))
	handler.Disk = readOnlyFakeDrive{}
	c.SetDE(0x005C)
	c.C = FunctionMakeFile
	if _, err := handler.Call(c, mem); err != nil {
		t.Fatalf("make read-only: %v", err)
	}
	if c.A != 0xFF {
		t.Fatalf("read-only make A=0x%02X", c.A)
	}
}

func writeFCB(mem cpu.Memory, addr uint16, name string, ext string) {
	mem.Write(addr, 0)
	for i := 0; i < 8; i++ {
		value := byte(' ')
		if i < len(name) {
			value = name[i]
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
}

func writeRenameFCB(mem cpu.Memory, addr uint16, name string, ext string) {
	for i := 0; i < 8; i++ {
		value := byte(' ')
		if i < len(name) {
			value = name[i]
		}
		mem.Write(addr+17+uint16(i), value)
	}
	for i := 0; i < 3; i++ {
		value := byte(' ')
		if i < len(ext) {
			value = ext[i]
		}
		mem.Write(addr+25+uint16(i), value)
	}
}

type fakeDrive struct {
	files map[string][]byte
}

func (d fakeDrive) List() ([]disk.Entry, error) {
	entries := make([]disk.Entry, 0, len(d.files))
	for name, data := range d.files {
		entries = append(entries, disk.Entry{Name: name, Size: int64(len(data))})
	}
	return entries, nil
}

func (d fakeDrive) ReadFile(name string) ([]byte, error) {
	data, ok := d.files[name]
	if !ok {
		return nil, errors.New("missing")
	}
	return append([]byte(nil), data...), nil
}

func (d fakeDrive) WriteFile(name string, data []byte) error {
	d.files[name] = append([]byte(nil), data...)
	return nil
}

func (d fakeDrive) DeleteFile(name string) error {
	if _, ok := d.files[name]; !ok {
		return errors.New("missing")
	}
	delete(d.files, name)
	return nil
}

func (d fakeDrive) RenameFile(oldName string, newName string) error {
	data, ok := d.files[oldName]
	if !ok {
		return errors.New("missing")
	}
	delete(d.files, oldName)
	d.files[newName] = data
	return nil
}

type readOnlyFakeDrive struct{}

func (readOnlyFakeDrive) List() ([]disk.Entry, error) { return nil, nil }

func (readOnlyFakeDrive) ReadFile(string) ([]byte, error) { return nil, errors.New("missing") }
