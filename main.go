package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const (
	//PORT web port
	PORT        = "80"
	ContentFile = ".ct"
	CommentFile = ".cm"
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
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("listening and serving on port :" + PORT)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template, _ := template.ParseFiles("./templates/index.html")

	template.Execute(w, nil)
}

// Handles sajak according to id
func ContentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	content, err := ReadContent(id)
	comments, err2 := ReadComments(id)
	if err != nil || err2 != nil {
		fmt.Println(err)
	}

	fmt.Println(comments)
	var page ContentPage
	page = ContentPage{content, comments}

	template, tmplerr := template.ParseFiles("./templates/content.html")
	fmt.Println(tmplerr)
	template.Execute(w, page)
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	r.ParseForm()

	c := Comment{r.Form["name"][0], r.Form["comment"][0]}

	if c.Name == "" {
		c.Name = "anon"
	}
	if c.CommentContent != "" {
		WriteComment(c, id)
	}
	http.Redirect(w, r, r.URL.Path, 302)
}

func ReadContent(fn string) (Content, error) {
	f, err := os.Open("./contents/" + fn + ContentFile)
	if err != nil {
		f.Close()
		return Content{}, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	s := strings.Join(text, "\n")

	title := strings.Split(s, "{seg}")[0]
	content := strings.Split(s, "{seg}")[1]
	verses := strings.Split(content, "{n}")

	sajak := Content{title, verses}

	f.Close()
	return (sajak), err
}

func WriteComment(c Comment, id string) {

	filename := "./comments/" + id + CommentFile
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		ioutil.WriteFile(filename, nil, 0600)
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	defer f.Close()
	text := ""

	info, _ := os.Stat(filename)

	if info.Size() > 0 {
		text += "{cm}\n"
	}

	text = text + c.Name + "\n{seg}\n" + c.CommentContent + "\n"
	fmt.Println(text)
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func ReadComments(id string) ([]Comment, error) {
	f, err := os.Open("./comments/" + id + CommentFile)
	if err != nil {
		f.Close()
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	info, _ := os.Stat("./comments/" + id + CommentFile)
	if info.Size() == 0 {
		return nil, nil
	}

	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	f.Close()
	s := strings.Join(text, "\n")

	rawComments := strings.Split(s, "{cm}")

	var comments []Comment

	for _, comm := range rawComments {
		spltComms := strings.Split(comm, "{seg}")
		comments = append(comments, Comment{spltComms[0], spltComms[1]})
	}

	return comments, nil
}

type ContentPage struct {
	Content  Content
	Comments []Comment
}

//Content is a poetry consists of Title and Verses
type Content struct {
	Title  string
	Verses []string
}

type Comment struct {
	Name           string
	CommentContent string
}
