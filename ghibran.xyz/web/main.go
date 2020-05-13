package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"ghibran.xyz/blogdata"
)

const (
	//PORT web port
	PORT = "80"

	// database settings
	DBUser     = "website"
	DBPassword = "!WebsitePW"
	DBIP       = "157.230.125.71"
	DBName     = "websiteDB"
)

var router *mux.Router
var database *sql.DB

func main() {
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPassword, DBIP, DBName)
	database, _ = sql.Open("mysql", connectionStr)

	dbError := database.Ping()
	if dbError != nil {
		log.Fatal(dbError)
	}
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
	router.HandleFunc("/favicon.ico", ServeFavicon)
	router.HandleFunc("/favicon.ico?v=2", ServeFavicon)

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
	router.HandleFunc("/admin/edit", AdminEditHandler).Methods("GET")
	router.HandleFunc("/admin/edit/{id}", PostEditHandler).Methods("GET")
	router.HandleFunc("/admin/edit/{id}", EditPostHandler).Methods("POST")
	router.HandleFunc("/admin/delete/{id}", DeletePostHandler).Methods("POST")

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/index.html")

	feed := blogdata.GetContentsCut(7, 100, 5, database)

	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, feed)

}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func ContentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	content, contenterr := blogdata.GetContent(id, database)
	comments, _ := blogdata.ReadComments(id, database)

	if contenterr != nil {
		fmt.Println(contenterr)
		return
	}

	var page blogdata.ContentPage

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
	if comment == "" || strings.TrimSpace(comment) == "" {
		http.Redirect(w, r, r.URL.Path, 302)
	}
	if name == "" {
		name = "Anon"
	}

	c := blogdata.Comment{Name: name, CommentContent: comment}
	blogdata.WriteComment(c, id, database)
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

	id := blogdata.WriteContent(c, database)

	http.Redirect(w, r, "/content/"+strconv.FormatInt(id, 10), 302)
}

func AdminEditHandler(w http.ResponseWriter, r *http.Request) {
	posts := blogdata.GetContentsCut(10, 100, 1, database)

	template, err := template.ParseFiles("./templates/editlist.html")
	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, posts)
}

func PostEditHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	content, err := blogdata.GetContent(id, database)
	if err != nil {
		log.Fatal(err)
	}

	template, err := template.ParseFiles("./templates/editpost.html")

	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, content)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	r.ParseForm()

	form := r.Form
	title := form["title"][0]
	text := form["content"][0]

	verses := strings.Split(text, "\n")

	content := blogdata.Content{Title: title, Verses: verses}

	blogdata.EditContent(id, &content, database)

	http.Redirect(w, r, "/content/"+id, 302)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	fmt.Println("delete ?" + id)
	blogdata.DeleteContent(id, database)

	blogdata.DeleteComments(id, database)

}
