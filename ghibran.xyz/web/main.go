package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"ghibran.xyz/blogdata"
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

	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/content/{id}", ContentHandler).Methods("GET")
	r.HandleFunc("/content/{id}", CommentHandler).Methods("POST")

	http.Handle("/", r)
	log.Println("listening and serving on port :" + PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template, _ := template.ParseFiles("/templates/index.html")

	template.Execute(w, nil)
}

// Handles sajak according to id
func ContentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	content, err := blogdata.ReadContent(id)
	comments, err2 := blogdata.ReadComments(id)
	if err != nil || err2 != nil {
		fmt.Println(err)
	}
	var page blogdata.ContentPage


	if(comments != nil){
		*comments = reverse(*comments)
	}

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

	c := blogdata.Comment{Name: r.Form["name"][0], CommentContent: r.Form["comment"][0]}

	if c.Name == "" {
		c.Name = "anon"
	}

	if c.CommentContent != "" {
		blogdata.WriteComment(c, id)
	}
	http.Redirect(w, r, r.URL.Path, 302)
}

func reverse(n []blogdata.Comment) []blogdata.Comment {
	for i := 0; i < len(n)/2; i++ {
		j := len(n) - i - 1
		n[i], n[j] = n[j], n[i]
	}
	return n
}
