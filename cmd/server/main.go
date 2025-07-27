package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mostlygeek/vibecities/db"
	"github.com/mostlygeek/vibecities/web"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:1337", "HTTP listen address")
	dbPath := flag.String("db", "vibecities.db", "Path to database file")
	flag.Parse()

	db, err := db.NewDBSqlite(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, ok := db.Get("/about"); !ok {
		db.Set("/about", "About Page", "VibeCities your place on the cool web.")
	}

	mcpSrv := web.NewMCPServer(db)
	srv := web.NewServer(db, mcpSrv)

	log.Println("Listening on " + *listen)
	http.ListenAndServe(*listen, srv)
}
