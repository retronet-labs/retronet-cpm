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
