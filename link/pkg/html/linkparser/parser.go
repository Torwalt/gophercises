package linkparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func GetLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("could not parse html: %v", err)
	}

	var links []Link
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				fmt.Println(a.Val)
				fmt.Println(a.Key)
				if a.Key == "href" {
					links = append(links, Link{Href: strings.TrimSpace(a.Val), Text: strings.TrimSpace(n.FirstChild.Data)})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(doc)

	return links, nil
}
