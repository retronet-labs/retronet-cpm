package disk

// Entry descrive un file visibile nel drive CP/M-like.
type Entry struct {
	Name string
	Size int64
}

// Drive e' il minimo contratto che il core CP/M usa per riferirsi a un disco.
type Drive interface {
	List() ([]Entry, error)
	ReadFile(name string) ([]byte, error)
}

// MutableDrive estende Drive con operazioni host mutanti. retronet-cpm lo usa
// solo quando la CLI abilita esplicitamente -write-disk.
type MutableDrive interface {
	Drive
	WriteFile(name string, data []byte) error
	DeleteFile(name string) error
	RenameFile(oldName string, newName string) error
}
