package db

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

type DBSqlite struct {
	sync.RWMutex
	db *sql.DB
}

func NewDBSqlite(path string) (*DBSqlite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create the records table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS records (
		path TEXT PRIMARY KEY,
		data TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DBSqlite{
		db: db,
	}, nil
}

func (d *DBSqlite) List() map[string]Record {
	d.RLock()
	defer d.RUnlock()

	rows, err := d.db.Query("SELECT path, data FROM records")
	if err != nil {
		return make(map[string]Record)
	}
	defer rows.Close()

	result := make(map[string]Record)
	for rows.Next() {
		var path, data string
		if err := rows.Scan(&path, &data); err != nil {
			continue
		}
		result[path] = Record{Data: data}
	}

	return result
}

func (d *DBSqlite) Set(path string, data string) error {
	d.Lock()
	defer d.Unlock()

	_, err := d.db.Exec("INSERT OR REPLACE INTO records (path, data) VALUES (?, ?)", path, data)
	return err
}

func (d *DBSqlite) Get(path string) (rec Record, ok bool) {
	d.RLock()
	defer d.RUnlock()

	var data string
	err := d.db.QueryRow("SELECT data FROM records WHERE path = ?", path).Scan(&data)
	if err != nil {
		return Record{}, false
	}

	return Record{Data: data}, true
}

func (d *DBSqlite) Close() error {
	return d.db.Close()
}
