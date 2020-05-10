package blogdata

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func WriteContent(c *Content) int {
	id := 0
	sfilename := "./contents/" + "stats"
	sf, err := os.OpenFile(sfilename, os.O_RDWR, 0644)
	if err != nil {
		sf, err = os.OpenFile(sfilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		id = 1
	} else {
		scanner := bufio.NewScanner(sf)
		scanner.Split(bufio.ScanLines)
		var txt string
		if scanner.Scan() {
			txt = scanner.Text()
		}
		id, err = strconv.Atoi(txt)

		if err != nil {
			id = 1
		}
	}
	id = id + 1
	fmt.Println(id)
	val := strconv.Itoa(id)
	fmt.Println(val)
	sf.WriteAt([]byte(val), 0)
	sf.Close()

	filename := "./contents/" + strconv.Itoa(id-1) + ContentFile
	fmt.Println(filename)
	cf, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	ctxt := c.Title + "\n{seg}\n"

	for _, s := range c.Verses {

		ctxt += s + "{n}\n"
	}

	cf.WriteAt([]byte(ctxt), 0)

	cf.Close()

	return id - 1
}

type ContentPage struct {
	Content  Content
	Comments []Comment
}
type Content struct {
	Title  string
	Verses []string
}
