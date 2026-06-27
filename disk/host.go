package disk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ErrInvalidName = fmt.Errorf("nome CP/M 8.3 non valido")
var ErrReadOnly = fmt.Errorf("drive host read-only")
var ErrFileTooLarge = fmt.Errorf("file CP/M oltre il limite configurato")
var ErrTooManyFiles = fmt.Errorf("drive CP/M oltre il numero massimo di file")

type HostDriveOptions struct {
	Writable    bool
	MaxFileSize int64
	MaxFiles    int
}

// HostDrive espone una directory host come drive A: read-only.
type HostDrive struct {
	root        string
	writable    bool
	maxFileSize int64
	maxFiles    int
}

func NewHostDrive(root string) (*HostDrive, error) {
	return NewHostDriveWithOptions(root, HostDriveOptions{})
}

func NewHostDriveWithOptions(root string, options HostDriveOptions) (*HostDrive, error) {
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
	return &HostDrive{
		root:        abs,
		writable:    options.Writable,
		maxFileSize: options.MaxFileSize,
		maxFiles:    options.MaxFiles,
	}, nil
}

func NewWritableHostDrive(root string) (*HostDrive, error) {
	return NewHostDriveWithOptions(root, HostDriveOptions{Writable: true})
}

func NewTemporaryHostDrive(prefix string, options HostDriveOptions) (*HostDrive, func() error, error) {
	if strings.TrimSpace(prefix) == "" {
		prefix = "retronet-cpm-"
	}
	root, err := os.MkdirTemp("", prefix)
	if err != nil {
		return nil, nil, err
	}
	drive, err := NewHostDriveWithOptions(root, options)
	if err != nil {
		_ = os.RemoveAll(root)
		return nil, nil, err
	}
	return drive, func() error { return os.RemoveAll(root) }, nil
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
	data, err := os.ReadFile(hostPath)
	if err != nil {
		return nil, err
	}
	if err := d.checkFileSize(int64(len(data))); err != nil {
		return nil, err
	}
	return data, nil
}

func (d *HostDrive) WriteFile(name string, data []byte) error {
	if !d.writable {
		return ErrReadOnly
	}
	if err := d.checkFileSize(int64(len(data))); err != nil {
		return err
	}
	normalized, err := NormalizeName(name)
	if err != nil {
		return err
	}
	if _, err := d.findHostPath(normalized); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := d.checkFileCountForNewFile(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return os.WriteFile(filepath.Join(d.root, normalized), data, 0o600)
}

func (d *HostDrive) DeleteFile(name string) error {
	if !d.writable {
		return ErrReadOnly
	}
	hostPath, err := d.findHostPath(name)
	if err != nil {
		return err
	}
	return os.Remove(hostPath)
}

func (d *HostDrive) RenameFile(oldName string, newName string) error {
	if !d.writable {
		return ErrReadOnly
	}
	oldPath, err := d.findHostPath(oldName)
	if err != nil {
		return err
	}
	normalized, err := NormalizeName(newName)
	if err != nil {
		return err
	}
	if _, err := d.findHostPath(normalized); err == nil {
		return fmt.Errorf("%w: %s", os.ErrExist, normalized)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.Rename(oldPath, filepath.Join(d.root, normalized))
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

func (d *HostDrive) checkFileSize(size int64) error {
	if d.maxFileSize > 0 && size > d.maxFileSize {
		return fmt.Errorf("%w: %d > %d", ErrFileTooLarge, size, d.maxFileSize)
	}
	return nil
}

func (d *HostDrive) checkFileCountForNewFile() error {
	if d.maxFiles <= 0 {
		return nil
	}
	entries, err := d.List()
	if err != nil {
		return err
	}
	if len(entries) >= d.maxFiles {
		return fmt.Errorf("%w: %d >= %d", ErrTooManyFiles, len(entries), d.maxFiles)
	}
	return nil
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
