package blogdata

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	CommentFile = ".cm"
)

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

	text = text + c.Name + " {seg} " + c.CommentContent + "\n"
	fmt.Println(text)
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func ReadComments(id string) (*[]Comment, error) {
	f, err := os.Open("./comments/" + id + CommentFile)
	
	if err != nil {
		fmt.Println(err)
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

	return &comments, nil
}

type Comment struct {
	Name           string
	CommentContent string
}
