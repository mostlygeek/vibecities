package main

import (
	"log"
	"net/http"

	"github.com/mostlygeek/mcpcities/db"
	"github.com/mostlygeek/mcpcities/web"
)

func main() {
	db := db.New("")
	db.Set("/about", "About Content")
	srv := web.NewServer(db)

	listen := "localhost:8111"
	log.Println("Listening on " + listen)
	http.ListenAndServe(listen, srv)
}
