package blogdata

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	ContentFile = ".ct"

	fromFile = false
)

func ReadContent(id string, database *sql.DB) (*Content, error) {
	queryCommand := fmt.Sprintf("SELECT title, text FROM contents WHERE id='%s'", id)

	rows, qerr := database.Query(queryCommand)
	defer rows.Close()
	var content Content
	if qerr != nil {
		return &content, qerr
	}

	if rows.Next() {
		var s string
		rows.Scan(&content.Title, &s)

		content.Verses = strings.Split(s, "{n}")
	} else {
		qerr = errors.New("404")
	}
	return &content, qerr
}

func WriteContent(c *Content, database *sql.DB) int64 {

	text := strings.Join(c.Verses, "{n}")

	queryCommand := "INSERT INTO contents (title , text) VALUES (?, ?)"

	prep, ierr := database.Prepare(queryCommand)

	if ierr != nil {
		log.Fatal(ierr)
	}
	defer prep.Close()

	res, reserr := prep.Exec(c.Title, text)

	if reserr != nil {
		log.Fatal(reserr)
		return -1
	}

	id, _ := res.LastInsertId()
	return id
}

func GetFeed(num int, database *sql.DB) *[]Content {
	contents := []Content{}

	queryCommand := fmt.Sprintf("SELECT id, title, text FROM contents ORDER BY id DESC LIMIT %d", num)

	query, err := database.Query(queryCommand)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	for query.Next() {
		var content Content
		s := ""
		query.Scan(&content.Id, &content.Title, &s)
		content.Verses = strings.Split(s, "{n}")[:3]
		content.Verses = append(content.Verses, ".......")
		contents = append(contents, content)
	}

	return &contents
}

type ContentPage struct {
	Content  Content
	Comments []Comment
}
type Content struct {
	Id     int
	Title  string
	Verses []string
}
