package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-8080/cpu"
	"github.com/retronet-labs/retronet-cpm/bdos"
	"github.com/retronet-labs/retronet-cpm/cpm"
	"github.com/retronet-labs/retronet-cpm/disk"
)

func TestSessionPromptCommandDrainAndSnapshot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.TXT"), []byte("CIAO"), 0o600); err != nil {
		t.Fatal(err)
	}
	drive, err := disk.NewHostDrive(root)
	if err != nil {
		t.Fatal(err)
	}
	sess, err := New(Config{Drive: drive})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := sess.Prompt(); err != nil {
		t.Fatal(err)
	}
	if err := sess.RunCommand("DIR"); err != nil {
		t.Fatal(err)
	}
	out, err := sess.DrainOutput()
	if err != nil {
		t.Fatal(err)
	}
	if text := string(out); !strings.Contains(text, "A>") || !strings.Contains(text, "README.TXT") {
		t.Fatalf("output=%q", text)
	}
	snap, err := sess.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if snap.Width != 80 || snap.Height != 24 {
		t.Fatalf("snapshot=%+v", snap)
	}
}

func TestSessionProgramInputUsesTerminalQueue(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "ECHO.COM"), echoProgram(), 0o600); err != nil {
		t.Fatal(err)
	}
	drive, err := disk.NewHostDrive(root)
	if err != nil {
		t.Fatal(err)
	}
	sess, err := New(Config{Drive: drive})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := sess.Input([]byte("Z")); err != nil {
		t.Fatal(err)
	}
	if err := sess.RunCommand("RUN ECHO"); err != nil {
		t.Fatal(err)
	}
	out, err := sess.DrainOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Z") || !strings.Contains(string(out), string(cpm.RunStoppedBDOSTerminate)) {
		t.Fatalf("output=%q", string(out))
	}
}

func echoProgram() []byte {
	return []byte{
		cpu.MVI(cpu.RegC), bdos.FunctionConsoleInput,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
		cpu.MVI(cpu.RegC), bdos.FunctionTerminate,
		cpu.CALL(), byte(cpm.BDOSVector), byte(cpm.BDOSVector >> 8),
	}
}
