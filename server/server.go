package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mostlygeek/mcpcities/db"
)

type Server struct {
	db  db.Store
	mux *http.ServeMux
}

func NewServer(db db.Store) *Server {
	s := &Server{
		db:  db,
		mux: http.NewServeMux(),
	}
	s.register()
	return s
}

func (s *Server) register() {
	s.mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			fmt.Fprintln(w, "Welcome to the index page!")
		} else {
			if rec, ok := s.db.Get(path); ok {
				io.WriteString(w, rec.Data)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
