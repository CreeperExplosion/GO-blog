package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"ghibran.xyz/blogdata"
)

const (
	//PORT web port
	PORT = "80"
)

var router *mux.Router

func main() {
	router = mux.NewRouter()

	ServeStatic()

	ServePage()

	http.Handle("/", router)
	log.Println("listening and serving on port :" + PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ServeStatic() {
	router.PathPrefix("/static/css").Handler(http.StripPrefix("/static/css",
		http.FileServer(http.Dir("./static/css"))))
	router.PathPrefix("/static/images").Handler(http.StripPrefix("/static/images",
		http.FileServer(http.Dir("./static/images"))))
}

func ServePage() {
	router.HandleFunc("/", IndexHandler).Methods("GET")
	router.HandleFunc("/content/{id}", ContentHandler).Methods("GET")
	router.HandleFunc("/content/{id}", CommentHandler).Methods("POST")
	router.HandleFunc("/admin", AdminHandler).Methods("GET")
	router.HandleFunc("/admin", AdminPostHandler).Methods("POST")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	} else {
		template.Execute(w, nil)
	}

}

func ContentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	content, contenterr := blogdata.ReadContent(id)
	comments, _ := blogdata.ReadComments(id)

	if contenterr != nil {
		fmt.Println(contenterr)
		return
	}

	var page blogdata.ContentPage

	*comments = blogdata.Reverse(*comments)

	page = blogdata.ContentPage{Content: *content, Comments: *comments}

	template, tmplerr := template.ParseFiles("./templates/content.html")

	if tmplerr != nil {
		log.Fatal(tmplerr)
	} else {
		template.Execute(w, page)
	}
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	r.ParseForm()

	name, comment := r.Form["name"][0], r.Form["comment"][0]

	if comment == "" {
		http.Redirect(w, r, r.URL.Path, 302)
	}

	if name == "" {
		name = "Anon"
	}

	c := blogdata.Comment{Name: name, CommentContent: comment}
	blogdata.WriteComment(c, id)
	http.Redirect(w, r, r.URL.Path, 302)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/admin.html")

	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
func AdminPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	form := r.Form
	content := form["content"][0]
	title := form["title"][0]

	fmt.Println(content)
	c := &blogdata.Content{Title: title, Verses: strings.Split(content, "\n")}

	id := blogdata.WriteContent(c)

	http.Redirect(w, r, "/content/"+strconv.Itoa(id), 302)
}
