package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"link/pkg/html/linkparser"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("could not run: %v", err)
	}
}

type XMLURL struct {
	Loc string `xml:"loc"`
}

type SiteMap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNs   string   `xml:"xmlns,attr"`
	URLs    []XMLURL `xml:"url"`
}

func run() error {
	urlP := flag.String("url", "https://www.calhoun.io/", "url to build a sitemap")
	flag.Parse()
	rURL := *urlP

	pURL, err := url.Parse(rURL)
	if err != nil {
		return fmt.Errorf("could not parse url: %v", err)
	}
	// stores unique links of the site
	linkMap := make(map[string]linkparser.Link)

	err = getLinksFromURL(rURL, pURL, linkMap)
	if err != nil {
		return fmt.Errorf("could not get links from url: %v", err)
	}

	allLinks := getKeys(linkMap)

	xmlUrls := []XMLURL{}
	for _, link := range allLinks {
		xmlUrls = append(xmlUrls, XMLURL{Loc: link})
	}

	sitemap := SiteMap{
		XMLNs: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  xmlUrls,
	}
	out, _ := xml.MarshalIndent(sitemap, " ", "  ")

	fmt.Println(xml.Header + string(out))

	return nil
}

// recursivly get links from an url by following all links and adding them to the map
func getLinksFromURL(u string, pURL *url.URL, linkMap map[string]linkparser.Link) error {
	// get html of url
	// sleep for a second to avoid overloading server
	time.Sleep(1 * time.Second)
	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("could not get content of url: %v", err)
	}
	defer resp.Body.Close()

	// get links from html
	links, err := linkparser.GetLinks(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse links: %v", err)
	}

	// clean links i.e.:
	// remove href duplicates & remove links to external sites
	newLinks := []string{}
	for _, link := range links {
		gl := link.Href

		pLink, err := url.Parse(link.Href)
		if err != nil {
			//fmt.Println(fmt.Errorf("could not parse link: %v", err))
			continue
		}

		// empty Host means local link
		if pLink.Host != "" && pLink.Host != pURL.Host {
			//fmt.Println(fmt.Sprintf("skipping link: %v\n as it leads to external site", link.Href))
			continue
		}

		// only work on schema of original url
		if pLink.Scheme != "" && pLink.Scheme != pURL.Scheme {
			continue
		}

		// ignore URL Fragments
		if pLink.Fragment != "" {
			continue
		}

		// build absolute link if local link
		if pLink.Host == "" {
			gl = fmt.Sprintf("%s%s", pURL.Host, gl)
		}
		if pLink.Scheme == "" {
			gl = fmt.Sprintf("%s://%s", pURL.Scheme, gl)
		}

		// check if link is already in map
		if _, ok := linkMap[gl]; !ok {
			linkMap[gl] = link
			fmt.Println(gl)
			newLinks = append(newLinks, gl)
		}
	}

	for _, link := range newLinks {
		err = getLinksFromURL(link, pURL, linkMap)
		if err != nil {
			return fmt.Errorf("could not get links from url: %v", err)
		}
	}

	return nil
}

func getKeys(m map[string]linkparser.Link) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys

}
