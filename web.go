package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func apiHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	fmt.Fprintf(w, "%s", returnJSON(db))
}

func leadersHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	fmt.Fprintf(w, "%s", returnLeadersJSON(db))
}

func serveWeb(db *DB) {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/api/get/data", func(w http.ResponseWriter, r *http.Request) {
		apiHandler(w, r, db)
	})
	http.HandleFunc("/api/get/leaders", func(w http.ResponseWriter, r *http.Request) {
		leadersHandler(w, r, db)
	})
	log.Println("listening on localhost:4000")
	http.ListenAndServe(":4000", nil)
}
