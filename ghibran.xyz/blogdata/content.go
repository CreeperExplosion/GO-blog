package blogdata

import (
	"bufio"
	"os"
	"strings"
)

const (
	ContentFile = ".ct"
)

func ReadContent(fn string) (*Content, error) {
	f, err := os.Open("./contents/" + fn + ContentFile)
	if err != nil {
		f.Close()
		return &Content{}, err
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
	return (&sajak), err
}

type ContentPage struct {
	Content  Content
	Comments []Comment
}
type Content struct {
	Title  string
	Verses []string
}
