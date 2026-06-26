package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
)

func TestAssemblerHelloBDOSEndToEnd(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Clean(".."))
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	asmRoot := filepath.Clean(filepath.Join(repoRoot, "..", "retronet-asm"))
	if _, err := os.Stat(filepath.Join(asmRoot, "go.mod")); err != nil {
		t.Skip("repo sibling retronet-asm non disponibile")
	}

	outPath := filepath.Join(t.TempDir(), "HELLO.COM")
	cmd := exec.Command("go", "run", "./cmd/retronet-asm", "build",
		filepath.Join(repoRoot, "examples", "hello-bdos.asm"),
		"-o", outPath,
	)
	cmd.Dir = asmRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("assembler: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "assemblato") {
		t.Fatalf("output assembler inatteso:\n%s", output)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("lettura COM: %v", err)
	}
	console := bdos.NewMemoryConsole(nil)
	m, err := cpm.NewMachine(cpm.Config{Console: console})
	if err != nil {
		t.Fatalf("NewMachine: %v", err)
	}
	if err := m.LoadCOM("HELLO.COM", data); err != nil {
		t.Fatalf("LoadCOM: %v", err)
	}
	result, err := m.Run(100)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Reason != cpm.RunStoppedBDOSTerminate || console.Output() != "HI" {
		t.Fatalf("result=%+v output=%q", result, console.Output())
	}
}
