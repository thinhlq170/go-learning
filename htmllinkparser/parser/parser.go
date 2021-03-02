package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link type
type Link struct {
	Href, Text string
}

//LinkParser return list Link struct
func LinkParser(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var listAtag []Link
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						listAtag = append(listAtag, Link{strings.TrimSpace(a.Val), strings.TrimSpace(n.LastChild.Data)})
					} else {
						break
					}
				}
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return listAtag, nil
}
