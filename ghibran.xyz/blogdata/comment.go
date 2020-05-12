package blogdata

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	CommentFile = ".cm"
)

func WriteComment(c Comment, id string, database *sql.DB) {

	queryCommand := "INSERT INTO comments (display_name , text, content_id) VALUES (?, ?, ?)"

	prep, ierr := database.Prepare(queryCommand)

	if ierr != nil {
		log.Fatal(ierr)
	}
	defer prep.Close()

	_, reserr := prep.Exec(c.Name, c.CommentContent, id)

	if reserr != nil {
		log.Fatal(reserr)

	}

}

func ReadComments(id string, database *sql.DB) (*[]Comment, error) {
	comments := []Comment{}

	queryCommand := fmt.Sprintf("SELECT display_name, text FROM comments WHERE content_id='%s' ORDER BY id DESC", id)

	query, err := database.Query(queryCommand)
	
	if err != nil {
		fmt.Println(err)
		return &comments, err
	}

	for query.Next() {
		var comment Comment

		query.Scan(&comment.Name, &comment.CommentContent)

		comments = append(comments, comment)
	}

	return &comments, err
}

type Comment struct {
	Name           string
	CommentContent string
}
