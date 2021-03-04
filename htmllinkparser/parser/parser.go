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

//recurse search text inside child node
func searchTextFunc(t *html.Node) string {
	text := ""
	for k := t.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.TextNode {
			text += " " + k.Data
		} else if k.Type == html.ElementNode { //skip tag and keep searching text inside the tag
			text += searchTextFunc(k)
		}
	}
	return strings.TrimSpace(text)
}

func linkParserFunc(n *html.Node) []Link {
	var listAtag []Link
	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					text := searchTextFunc(n)
					listAtag = append(listAtag, Link{strings.TrimSpace(a.Val), text})
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		listAtag = append(listAtag, linkParserFunc(c)...)
	}
	return listAtag
}

//LinkParser return list Link struct
func LinkParser(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	//listAtag := linkParserFunc(doc)
	// var listAtag []Link
	// var f func(*html.Node)

	// f = func(n *html.Node) {
	// 	if n.Type == html.ElementNode {
	// 		if n.Data == "a" {
	// 			for _, a := range n.Attr {
	// 				if a.Key == "href" {
	// 					text := searchTextFunc(n)
	// 					listAtag = append(listAtag, Link{strings.TrimSpace(a.Val), text})
	// 				} else {
	// 					break
	// 				}
	// 			}
	// 		}

	// 	}
	// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 		f(c)
	// 	}
	// }
	//f(doc)
	return linkParserFunc(doc), nil
}
