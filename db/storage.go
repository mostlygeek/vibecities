package db

type Store interface {
	List() map[string]Record
	Set(path string, data string) error
	Get(path string) (rec Record, ok bool)
	Delete(path string) error
}

type Record struct {
	Data string
}
