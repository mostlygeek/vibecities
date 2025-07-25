package db

import (
	"sync"
)

type Store interface {
	List() map[string]Record
	Set(path string, data string) error
	Get(path string) (rec Record, ok bool)
}

type Record struct {
	Data string
}

// eventually make this an sqliteDB
type DB struct {
	sync.RWMutex
	data map[string]Record
}

func New(dbpath string) *DB {
	return &DB{
		data: make(map[string]Record),
	}
}

func (d *DB) List() map[string]Record {
	d.RLock()
	defer d.RUnlock()

	return d.data
}

func (d *DB) Set(path string, data string) error {
	d.Lock()
	defer d.Unlock()
	d.data[path] = Record{Data: data}
	return nil
}

func (d *DB) Get(path string) (rec Record, ok bool) {
	d.RLock()
	defer d.RUnlock()

	rec, ok = d.data[path]
	return
}
