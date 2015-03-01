package main

import (
	"fmt"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", returnJSON())
}

func serveWeb() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/api/get/data", apiHandler)
	http.ListenAndServe(":4000", nil)
}
