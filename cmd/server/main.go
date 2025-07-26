package main

import (
	"log"
	"net/http"

	"github.com/mostlygeek/vibecities/db"
	"github.com/mostlygeek/vibecities/web"
)

func main() {
	db, err := db.NewDBSqlite("vibecities.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, ok := db.Get("/about"); !ok {
		db.Set("/about", "About Page", "VibeCities your place on the cool web.")
	}

	mcpSrv := web.NewMCPServer(db)
	srv := web.NewServer(db, mcpSrv)

	listen := "localhost:1337"
	log.Println("Listening on " + listen)
	http.ListenAndServe(listen, srv)
}
