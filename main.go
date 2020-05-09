package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	//PORT web port
	PORT = "80"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/css").Handler(http.StripPrefix("/static/css",
		http.FileServer(http.Dir("./static/css"))))
	r.PathPrefix("/static/images").Handler(http.StripPrefix("/static/images",
		http.FileServer(http.Dir("./static/images"))))

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/content", contentHandler).Methods("GET")

	http.Handle("/", r)
	err := http.ListenAndServe(":"+PORT, nil)

	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("listening and serving on port :" + PORT)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template, _ := template.ParseFiles("./templates/index.html")

	template.Execute(w, nil)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	template, _ := template.ParseFiles("./templates/content.html")

	template.Execute(w, nil)
}

//Sajak content of blog
type Sajak struct {
	Title, Content string
}
