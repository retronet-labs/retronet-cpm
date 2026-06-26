// Package bdos implementa il subset console del BDOS CP/M-like.
package bdos

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/disk"
)

const (
	FunctionTerminate        byte   = 0
	FunctionConsoleInput     byte   = 1
	FunctionConsoleOutput    byte   = 2
	FunctionDirectConsoleIO  byte   = 6
	FunctionPrintString      byte   = 9
	FunctionReadConsoleLine  byte   = 10
	FunctionConsoleStatus    byte   = 11
	FunctionReturnVersion    byte   = 12
	FunctionOpenFile         byte   = 15
	FunctionCloseFile        byte   = 16
	FunctionReadSequential   byte   = 20
	FunctionSetDMA           byte   = 26
	DefaultVersion           byte   = 0x22
	DefaultDMA               uint16 = 0x0080
	RecordSize               int    = 128
	directConsoleInputMarker byte   = 0xFF
)

var (
	ErrNilCPU              = errors.New("cpu 8080 non inizializzata")
	ErrNilMemory           = errors.New("memoria BDOS non inizializzata")
	ErrNilConsole          = errors.New("console BDOS non inizializzata")
	ErrNilDisk             = errors.New("drive BDOS non inizializzato")
	ErrUnsupportedFunction = errors.New("funzione BDOS non supportata")
	ErrUnterminatedString  = errors.New("stringa BDOS non terminata")
)

// Console e' il terminale logico usato dalle funzioni BDOS console.
type Console interface {
	ReadByte() (byte, error)
	WriteByte(value byte) error
	Status() bool
}

// CallResult descrive l'effetto ad alto livello di una chiamata BDOS.
type CallResult struct {
	Function  byte
	Terminate bool
}

// UnsupportedFunctionError conserva il numero della funzione BDOS sconosciuta.
type UnsupportedFunctionError struct {
	Function byte
}

func (e *UnsupportedFunctionError) Error() string {
	return fmt.Sprintf("%v: C=%d", ErrUnsupportedFunction, e.Function)
}

func (e *UnsupportedFunctionError) Unwrap() error { return ErrUnsupportedFunction }

// Handler esegue le funzioni BDOS sullo stato CPU/memoria passato dal core CP/M.
type Handler struct {
	Console Console
	Disk    disk.Drive
	Version byte
	DMA     uint16

	openFiles map[uint16]*openFile
}

type openFile struct {
	name string
	data []byte
}

// NewHandler crea un handler BDOS. Se console e' nil viene usata una console
// vuota che scarta l'output.
func NewHandler(console Console) *Handler {
	if console == nil {
		console = NewMemoryConsole(nil)
	}
	h := &Handler{Console: console, Version: DefaultVersion}
	h.Reset()
	return h
}

// Reset ripristina lo stato volatile BDOS: DMA e file aperti.
func (h *Handler) Reset() {
	if h == nil {
		return
	}
	h.DMA = DefaultDMA
	h.openFiles = make(map[uint16]*openFile)
}

// Call esegue la funzione indicata dal registro C.
func (h *Handler) Call(c *cpu.CPU8080, mem cpu.Memory) (CallResult, error) {
	if c == nil {
		return CallResult{}, ErrNilCPU
	}
	if mem == nil {
		return CallResult{}, ErrNilMemory
	}
	if h == nil || h.Console == nil {
		return CallResult{}, ErrNilConsole
	}
	result := CallResult{Function: c.C}
	switch c.C {
	case FunctionTerminate:
		result.Terminate = true
	case FunctionConsoleInput:
		value, err := h.Console.ReadByte()
		if err != nil {
			return result, err
		}
		c.A = value
		c.L = value
		if err := h.Console.WriteByte(value); err != nil {
			return result, err
		}
	case FunctionConsoleOutput:
		return result, h.Console.WriteByte(c.E)
	case FunctionDirectConsoleIO:
		if c.E == directConsoleInputMarker {
			if !h.Console.Status() {
				c.A = 0
				c.L = 0
				return result, nil
			}
			value, err := h.Console.ReadByte()
			if err != nil {
				return result, err
			}
			c.A = value
			c.L = value
			return result, nil
		}
		return result, h.Console.WriteByte(c.E)
	case FunctionPrintString:
		return result, h.printDollarString(c, mem)
	case FunctionReadConsoleLine:
		return result, h.readConsoleLine(c, mem)
	case FunctionConsoleStatus:
		value := byte(0)
		if h.Console.Status() {
			value = 0xFF
		}
		c.A = value
		c.L = value
	case FunctionReturnVersion:
		version := h.Version
		if version == 0 {
			version = DefaultVersion
		}
		c.A = version
		c.H = 0
		c.L = version
	case FunctionOpenFile:
		return result, h.openFile(c, mem)
	case FunctionCloseFile:
		return result, h.closeFile(c)
	case FunctionReadSequential:
		return result, h.readSequential(c, mem)
	case FunctionSetDMA:
		h.DMA = c.DE()
	default:
		return result, &UnsupportedFunctionError{Function: c.C}
	}
	return result, nil
}

func (h *Handler) printDollarString(c *cpu.CPU8080, mem cpu.Memory) error {
	start := c.DE()
	for offset := 0; offset <= cpu.AddressSpaceSize; offset++ {
		addr := start + uint16(offset)
		value := mem.Read(addr)
		if value == '$' {
			return nil
		}
		if err := h.Console.WriteByte(value); err != nil {
			return err
		}
	}
	return ErrUnterminatedString
}

func (h *Handler) readConsoleLine(c *cpu.CPU8080, mem cpu.Memory) error {
	addr := c.DE()
	max := mem.Read(addr)
	countAddr := addr + 1
	dataAddr := addr + 2
	count := byte(0)
	for count < max {
		value, err := h.Console.ReadByte()
		if err != nil {
			return err
		}
		if value == '\n' {
			value = '\r'
		}
		if value == '\r' {
			break
		}
		mem.Write(dataAddr+uint16(count), value)
		count++
		if err := h.Console.WriteByte(value); err != nil {
			return err
		}
	}
	mem.Write(countAddr, count)
	return nil
}

func (h *Handler) openFile(c *cpu.CPU8080, mem cpu.Memory) error {
	if h.Disk == nil {
		return ErrNilDisk
	}
	addr := c.DE()
	name, err := fcbName(mem, addr)
	if err != nil {
		setReturnByte(c, 0xFF)
		return nil
	}
	data, err := h.Disk.ReadFile(name)
	if err != nil {
		setReturnByte(c, 0xFF)
		return nil
	}
	h.openFiles[addr] = &openFile{name: name, data: data}
	mem.Write(addr+32, 0)
	setReturnByte(c, 0)
	return nil
}

func (h *Handler) closeFile(c *cpu.CPU8080) error {
	addr := c.DE()
	if _, ok := h.openFiles[addr]; !ok {
		setReturnByte(c, 0xFF)
		return nil
	}
	delete(h.openFiles, addr)
	setReturnByte(c, 0)
	return nil
}

func (h *Handler) readSequential(c *cpu.CPU8080, mem cpu.Memory) error {
	addr := c.DE()
	file := h.openFiles[addr]
	if file == nil {
		setReturnByte(c, 0xFF)
		return nil
	}
	record := int(mem.Read(addr + 32))
	offset := record * RecordSize
	if offset >= len(file.data) {
		setReturnByte(c, 1)
		return nil
	}
	for i := 0; i < RecordSize; i++ {
		value := byte(0x1A)
		if offset+i < len(file.data) {
			value = file.data[offset+i]
		}
		mem.Write(h.DMA+uint16(i), value)
	}
	mem.Write(addr+32, byte(record+1))
	setReturnByte(c, 0)
	return nil
}

func fcbName(mem cpu.Memory, addr uint16) (string, error) {
	name := fcbPart(mem, addr+1, 8)
	ext := fcbPart(mem, addr+9, 3)
	if ext != "" {
		name += "." + ext
	}
	return disk.NormalizeName(name)
}

func fcbPart(mem cpu.Memory, addr uint16, size int) string {
	var b strings.Builder
	for i := 0; i < size; i++ {
		value := mem.Read(addr+uint16(i)) & 0x7F
		if value == ' ' || value == 0 {
			continue
		}
		b.WriteByte(value)
	}
	return b.String()
}

func setReturnByte(c *cpu.CPU8080, value byte) {
	c.A = value
	c.L = value
}

// MemoryConsole e' una console in memoria utile per test, conformance e CLI non
// interattiva.
type MemoryConsole struct {
	mu     sync.Mutex
	input  []byte
	output strings.Builder
}

func NewMemoryConsole(input []byte) *MemoryConsole {
	return &MemoryConsole{input: append([]byte(nil), input...)}
}

func (c *MemoryConsole) ReadByte() (byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.input) == 0 {
		return 0, io.EOF
	}
	value := c.input[0]
	c.input = c.input[1:]
	return value, nil
}

func (c *MemoryConsole) WriteByte(value byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.output.WriteByte(value)
}

func (c *MemoryConsole) Status() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.input) > 0
}

func (c *MemoryConsole) Output() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.output.String()
}

func (c *MemoryConsole) QueueInput(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.input = append(c.input, data...)
}

// StreamConsole collega BDOS a reader/writer host.
type StreamConsole struct {
	reader *bufio.Reader
	writer io.Writer
}

func NewStreamConsole(input io.Reader, output io.Writer) *StreamConsole {
	if input == nil {
		input = strings.NewReader("")
	}
	if output == nil {
		output = io.Discard
	}
	return &StreamConsole{reader: bufio.NewReader(input), writer: output}
}

func (c *StreamConsole) ReadByte() (byte, error) {
	value, err := c.reader.ReadByte()
	return value, err
}

func (c *StreamConsole) WriteByte(value byte) error {
	_, err := c.writer.Write([]byte{value})
	return err
}

func (c *StreamConsole) Status() bool {
	return c.reader.Buffered() > 0
}
