package shell

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
)

func TestShellCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "HELLO.TXT"), []byte("testo"), 0o600); err != nil {
		t.Fatalf("write text: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "HI.COM"), helloProgram(), 0o600); err != nil {
		t.Fatalf("write com: %v", err)
	}
	drive, err := disk.NewHostDrive(root)
	if err != nil {
		t.Fatalf("drive: %v", err)
	}
	var out bytes.Buffer
	sh, err := New(Config{
		Drive:  drive,
		Input:  strings.NewReader("DIR\nTYPE HELLO.TXT\nRUN HI ARG.TXT\nEXIT\n"),
		Output: &out,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := sh.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	text := out.String()
	for _, want := range []string{"A>", "HELLO.TXT", "HI.COM", "testo", "HI", string(cpm.RunStoppedBDOSTerminate)} {
		if !strings.Contains(text, want) {
			t.Fatalf("output missing %q:\n%s", want, text)
		}
	}
}

func TestShellProgramInputKeepsFollowingCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "ECHO.COM"), echoInputProgram(), 0o600); err != nil {
		t.Fatalf("write com: %v", err)
	}
	drive, err := disk.NewHostDrive(root)
	if err != nil {
		t.Fatalf("drive: %v", err)
	}
	var out bytes.Buffer
	sh, err := New(Config{
		Drive:  drive,
		Input:  strings.NewReader("RUN ECHO\nZEXIT\n"),
		Output: &out,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := sh.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "Z") || strings.Contains(text, "comando sconosciuto") {
		t.Fatalf("output:\n%s", text)
	}
}

func helloProgram() []byte {
	return []byte{
		cpu.LXI(cpu.PairDE), 0x0D, 0x01,
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		'H', 'I', '$',
	}
}

func echoInputProgram() []byte {
	return []byte{
		cpu.MVI(cpu.RegC), bdos.FunctionConsoleInput,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
}
