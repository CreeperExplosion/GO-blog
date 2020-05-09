package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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
	r.HandleFunc("/content/{id}", contentHandler).Methods("GET")

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
	params := mux.Vars(r)
	id := params["id"]

	content, err := ReadFile(id)
	if err != nil {
		println(err)
	}

	template, _ := template.ParseFiles("./templates/content.html")
	template.Execute(w, content)
}

func ReadFile(fn string) (Sajak, error) {
	f, err := os.Open("./content/" + fn + ".ct")
	if err != nil {
		f.Close()
		return Sajak{}, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	s := strings.Join(text, "")

	title := strings.Split(s, "~")[0]
	content := strings.Split(s, "~")[1]
	verses := strings.Split(content, "{n}")

	sajak := Sajak{title, verses}

	f.Close()
	return (sajak), err
}

type Sajak struct {
	Title  string
	Verses []string
}
