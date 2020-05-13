package blogdata

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

const (
	ContentFile = ".ct"

	fromFile = false
)

func GetContent(id string, database *sql.DB) (*Content, error) {
	queryCommand := `SELECT id, title, text FROM contents WHERE id=?`

	rows, qerr := database.Query(queryCommand, id)
	defer rows.Close()
	var content Content
	if qerr != nil {
		return &content, qerr
	}

	if rows.Next() {
		var s string
		rows.Scan(&content.Id, &content.Title, &s)

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

func GetContentsCut(num int, charlen int, versenum int, database *sql.DB) *[]Content {
	contents := []Content{}

	queryCommand := `SELECT id, title, text FROM contents ORDER BY id DESC LIMIT ?`

	query, err := database.Query(queryCommand, num)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	for query.Next() {
		var content Content
		s := ""
		query.Scan(&content.Id, &content.Title, &s)

		verses := strings.Split(s, "{n}")
		verselen := len(verses)

		if verselen > versenum {
			verses = verses[:versenum]
		}

		for _, verse := range verses {
			a := []rune(verse)

			leng := len(verse)
			if leng > charlen {
				leng = charlen
				content.Verses = append(content.Verses, string(a[0:leng])+"....")
			} else {
				content.Verses = append(content.Verses, verse)
			}
		}

		content.Verses = append(content.Verses, ".......")
		contents = append(contents, content)
	}

	return &contents
}

func DeleteContent(id string, database *sql.DB) {
	queryCommand := `DELETE FROM contents WHERE id=?`

	_, err := database.Exec(queryCommand, id)

	if err != nil {
		log.Fatal(err)
	}
}

func EditContent(id string, content *Content, database *sql.DB) {

	s := strings.Join(content.Verses, "{n}")
	queryCommand := `UPDATE contents SET title=?, text=? WHERE id=?`
	_, err := database.Exec(queryCommand, content.Title, s, id)

	if err != nil {
		log.Fatal(err)
	}

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
