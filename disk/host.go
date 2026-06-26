package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ErrInvalidName = fmt.Errorf("nome CP/M 8.3 non valido")

// HostDrive espone una directory host come drive A: read-only.
type HostDrive struct {
	root string
}

func NewHostDrive(root string) (*HostDrive, error) {
	if strings.TrimSpace(root) == "" {
		root = "."
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s non e' una directory", abs)
	}
	return &HostDrive{root: abs}, nil
}

func (d *HostDrive) Root() string { return d.root }

func (d *HostDrive) List() ([]Entry, error) {
	entries, err := os.ReadDir(d.root)
	if err != nil {
		return nil, err
	}
	out := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name, err := NormalizeName(entry.Name())
		if err != nil {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, Entry{Name: name, Size: info.Size()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (d *HostDrive) ReadFile(name string) ([]byte, error) {
	hostPath, err := d.findHostPath(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(hostPath)
}

func (d *HostDrive) findHostPath(name string) (string, error) {
	normalized, err := NormalizeName(name)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(d.root)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		entryName, err := NormalizeName(entry.Name())
		if err == nil && entryName == normalized {
			return filepath.Join(d.root, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("%w: %s", os.ErrNotExist, normalized)
}

// NormalizeName converte un nome host o CP/M in forma 8.3 maiuscola.
func NormalizeName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToUpper(value), "A:") {
		value = value[2:]
	}
	if value == "" || strings.ContainsAny(value, `/\`) || strings.Contains(value, ":") || strings.Contains(value, "..") {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, value)
	}
	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, value)
	}
	name := strings.ToUpper(parts[0])
	ext := ""
	if len(parts) == 2 {
		ext = strings.ToUpper(parts[1])
	}
	if len(name) == 0 || len(name) > 8 || len(ext) > 3 {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, value)
	}
	if !validPart(name) || (ext != "" && !validPart(ext)) {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, value)
	}
	if ext == "" {
		return name, nil
	}
	return name + "." + ext, nil
}

func validPart(value string) bool {
	for _, r := range value {
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		switch r {
		case '_', '$', '~', '-':
			continue
		default:
			return false
		}
	}
	return true
}
