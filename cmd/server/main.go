package main

import (
	"log"
	"net/http"

	"github.com/mostlygeek/vibecities/db"
	"github.com/mostlygeek/vibecities/web"
)

func main() {
	db, err := db.NewDBSqlite("mcpcities.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Set("/about", "About Content")

	mcpSrv := web.NewMCPServer(db)
	srv := web.NewServer(db, mcpSrv)

	listen := "localhost:8111"
	log.Println("Listening on " + listen)
	http.ListenAndServe(listen, srv)
}
