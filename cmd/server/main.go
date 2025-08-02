package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mostlygeek/vibecities/db"
	"github.com/mostlygeek/vibecities/web"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:1337", "HTTP listen address")
	dbPath := flag.String("db", "vibecities.db", "Path to database file")
	enableLoadPath := flag.Bool("enable-load-path", false, "Enable the loadPath tool (security risk)")

	// Custom usage message (shown with -h/--help)
	flag.Usage = func() {
		fmt.Println("Usage: server [options]")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println()
	}

	// If -help or -h supplied, flag.Parse will call usage and exit(2) for us
	flag.Parse()

	db, err := db.NewDBSqlite(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, ok := db.Get("/about"); !ok {
		db.Set("/about", "About Page", "VibeCities your place on the cool web.")
	}

	mcpSrv := web.NewMCPServer(db, *enableLoadPath)
	srv := web.NewServer(db, mcpSrv)

	log.Println("Listening on " + *listen)
	if err := http.ListenAndServe(*listen, srv); err != nil {
		log.Fatal(err)
	}
}
