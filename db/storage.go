package db

import (
	"database/sql"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Store interface {
	List() map[string]Record
	Set(path, title, data string) error
	Get(path string) (rec Record, ok bool)
	Delete(path string) error
}

type Record struct {
	ID      int64     `json:"id"`
	Path    string    `json:"path"`
	Title   string    `json:"title"`
	Data    string    `json:"data"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
type DBSqlite struct {
	sync.RWMutex
	db *sql.DB
}

func NewDBSqlite(path string) (*DBSqlite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY,
		path TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		data TEXT NOT NULL,
		created TIMESTAMP,
		updated TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DBSqlite{db: db}, nil
}

func (d *DBSqlite) List() map[string]Record {
	d.RLock()
	defer d.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, path, title, data, created, updated
		FROM records`)
	if err != nil {
		return make(map[string]Record)
	}
	defer rows.Close()

	result := make(map[string]Record)
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Path, &r.Title,
			&r.Data, &r.Created, &r.Updated); err != nil {
			continue
		}
		result[r.Path] = r
	}
	return result
}

func (d *DBSqlite) Set(path, title, data string) error {
	d.Lock()
	defer d.Unlock()

	now := time.Now().UTC()
	_, err := d.db.Exec(`
		INSERT INTO records (path, title, data, created, updated)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title=excluded.title,
			data=excluded.data,
			updated=excluded.updated`,
		path, title, data, now, now)
	return err
}

func (d *DBSqlite) Get(path string) (rec Record, ok bool) {
	d.RLock()
	defer d.RUnlock()

	var r Record
	err := d.db.QueryRow(`
		SELECT id, path, title, data, created, updated
		FROM records WHERE path = ?`, path).
		Scan(&r.ID, &r.Path, &r.Title, &r.Data, &r.Created, &r.Updated)
	if err != nil {
		return Record{}, false
	}
	return r, true
}

func (d *DBSqlite) Delete(path string) error {
	d.Lock()
	defer d.Unlock()

	_, err := d.db.Exec(`DELETE FROM records WHERE path = ?`, path)
	return err
}

func (d *DBSqlite) Close() error {
	return d.db.Close()
}
