package main

import (
	"log"
	"net/http"
)

func main() {
	webserver := http.FileServer(http.Dir("./static"))
	http.Handle("/", webserver)

	// fileserver := http.FileServer(http.Dir("./images/logo.png"))

	// http.Handle("images/", fileserver)

	log.Println("Listening on :80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
