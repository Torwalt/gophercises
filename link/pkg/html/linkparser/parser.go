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
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			href := getLinkFromA(n)
			text := getTextFromA(n)
			links = append(links, Link{Href: href, Text: text})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	return links, nil
}

func getLinkFromA(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "href" {
			return strings.TrimSpace(a.Val)
		}
	}
	return ""
}

func getTextFromA(n *html.Node) string {
	var allText []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			allText = append(allText, getTextFromA(c))
		}

		if c.Type == html.TextNode {
			allText = append(allText, strings.TrimSpace(c.Data))
		}
	}

	return strings.TrimSpace(strings.Join(allText, " "))
}
