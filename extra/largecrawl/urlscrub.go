package main

import (
	"html"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/miku/parallel"
)

// url.Parse:                               1369042
// url.Parse, Trim                         91672978
// url.Parse, Trim, Unescape               91629399
// url.Parse, Trim, Unescape, ReplaceAll   91663562

// SO: 91,663,562 URLs it is

func Scrub(p []byte) ([]byte, error) {
	s := string(p)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSpace(s)
	s = html.UnescapeString(s)
	_, err := url.Parse(s)
	if err != nil {
		return nil, nil
	}
	return p, nil
}

func main() {
	p := parallel.NewProcessor(os.Stdin, os.Stdout, Scrub)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}