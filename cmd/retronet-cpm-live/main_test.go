package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
)

func TestRunScriptedLiveShell(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.TXT"), []byte("testo live"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "HI.COM"), liveHelloProgram(), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"-disk", root,
		"-width", "80",
		"-height", "20",
		"-script", "DIR\rTYPE README.TXT\rRUN HI\rEXIT\r",
	}, strings.NewReader(""), &stdout, &stderr, false)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	output := stdout.String()
	for _, want := range []string{"RetroNet CP/M Live", "A>", "README.TXT", "testo live", "HI", string(cpm.RunStoppedBDOSTerminate)} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q in stdout=%s", want, output)
		}
	}
}

func TestInvalidALU(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-alu", "bad", "-script", "EXIT\r"}, strings.NewReader(""), &stdout, &stderr, false)
	if code != 2 || !strings.Contains(stderr.String(), "errore ALU") {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
}

func liveHelloProgram() []byte {
	return []byte{
		cpu.LXI(cpu.PairDE), 0x0D, 0x01,
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		'H', 'I', '$',
	}
}
