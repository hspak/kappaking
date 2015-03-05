package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func apiHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Fprintf(w, "%s", returnJSON(db))
}

func serveWeb(db *sql.DB) {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/api/get/data", func(w http.ResponseWriter, r *http.Request) {
		apiHandler(w, r, db)
	})
	log.Println("listening on localhost:4000")
	http.ListenAndServe(":4000", nil)
}
