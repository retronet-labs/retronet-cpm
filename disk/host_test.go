package disk

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeName(t *testing.T) {
	tests := map[string]string{
		"hello.com": "HELLO.COM",
		"A:read.me": "READ.ME",
		"FOO":       "FOO",
	}
	for input, want := range tests {
		got, err := NormalizeName(input)
		if err != nil {
			t.Fatalf("NormalizeName(%q): %v", input, err)
		}
		if got != want {
			t.Fatalf("NormalizeName(%q)=%q want %q", input, got, want)
		}
	}
	for _, input := range []string{"../X.COM", "TOO-LONG-NAME.COM", "A:B:C", "X.YYYY"} {
		if _, err := NormalizeName(input); !errors.Is(err, ErrInvalidName) {
			t.Fatalf("NormalizeName(%q) err=%v", input, err)
		}
	}
}

func TestHostDriveListAndRead(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "hello.com"), []byte{1, 2, 3}, 0o600); err != nil {
		t.Fatalf("write hello: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("test"), 0o600); err != nil {
		t.Fatalf("write notes: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "longfilename.txt"), []byte("skip"), 0o600); err != nil {
		t.Fatalf("write long: %v", err)
	}

	drive, err := NewHostDrive(root)
	if err != nil {
		t.Fatalf("NewHostDrive: %v", err)
	}
	entries, err := drive.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 || entries[0].Name != "HELLO.COM" || entries[1].Name != "NOTES.TXT" {
		t.Fatalf("entries=%+v", entries)
	}
	data, err := drive.ReadFile("A:HELLO.COM")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != string([]byte{1, 2, 3}) {
		t.Fatalf("data=%v", data)
	}
	if _, err := drive.ReadFile("..\\SECRET.COM"); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("path traversal err=%v", err)
	}
}

func TestHostDriveWriteRequiresExplicitWritableDrive(t *testing.T) {
	root := t.TempDir()
	drive, err := NewHostDrive(root)
	if err != nil {
		t.Fatalf("NewHostDrive: %v", err)
	}
	if err := drive.WriteFile("OUT.TXT", []byte("no")); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("read-only write err=%v", err)
	}

	writable, err := NewWritableHostDrive(root)
	if err != nil {
		t.Fatalf("NewWritableHostDrive: %v", err)
	}
	if err := writable.WriteFile("OUT.TXT", []byte("ok")); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	data, err := writable.ReadFile("out.txt")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "ok" {
		t.Fatalf("data=%q", data)
	}
	if err := writable.RenameFile("OUT.TXT", "NEW.TXT"); err != nil {
		t.Fatalf("RenameFile: %v", err)
	}
	if _, err := writable.ReadFile("OUT.TXT"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("old read err=%v", err)
	}
	if err := writable.DeleteFile("NEW.TXT"); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}
	if _, err := writable.ReadFile("NEW.TXT"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("deleted read err=%v", err)
	}
}
