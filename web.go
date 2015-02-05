package main

import "net/http"

func serveWeb() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.ListenAndServe(":3000", nil)
}
