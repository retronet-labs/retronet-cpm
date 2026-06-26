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

func TestRunConformance(t *testing.T) {
	var out, stderr bytes.Buffer
	code := run([]string{"-conformance"}, strings.NewReader(""), &out, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s stdout=%s", code, stderr.String(), out.String())
	}
	if !strings.Contains(out.String(), "conformance passed=") {
		t.Fatalf("stdout=%s", out.String())
	}
}

func TestRunProgramFromDisk(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "HI.COM"), cliHelloProgram(), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	var out, stderr bytes.Buffer
	code := run([]string{"-disk", root, "-run", "HI", "-alu", "native"}, strings.NewReader(""), &out, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s stdout=%s", code, stderr.String(), out.String())
	}
	for _, want := range []string{"HI", "program=HI.COM", string(cpm.RunStoppedBDOSTerminate)} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q in stdout=%s", want, out.String())
		}
	}
}

func TestInvalidALUEnv(t *testing.T) {
	t.Setenv("RETRONET_CPM_ALU", "bad")
	var out, stderr bytes.Buffer
	code := run([]string{"-conformance"}, strings.NewReader(""), &out, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "errore ALU") {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
}

func cliHelloProgram() []byte {
	return []byte{
		cpu.LXI(cpu.PairDE), 0x0D, 0x01,
		cpu.MVI(cpu.RegC), bdos.FunctionPrintString,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		'H', 'I', '$',
	}
}
