package builder

import (
	"encoding/xml"
	"io"
	"strings"

	"example.com/main/htmllinkparser/parser"
)

type urlBuilder struct {
	XMLName xml.Name `xml:"url"`
	LOC     string   `xml:"loc"`
}

//LinkBuilder gets list Link objects of the page and returns io.Writer and error
func LinkBuilder(page string, linksInPage []parser.Link, writtenFile io.Writer) error {
	var listLink []urlBuilder
	var errResult error
	if len(linksInPage) > 0 {
		enc := xml.NewEncoder(writtenFile)
		for _, link := range linksInPage {
			if strings.Contains(link.Href, page) {
				listLink = append(listLink, urlBuilder{LOC: link.Href})
			} else if link.Href[0] == '/' {
				href := page + link.Href[1:]
				listLink = append(listLink, urlBuilder{LOC: href})
			}
		}
		enc.Indent(" ", "	")
		if err := enc.Encode(listLink); err != nil {
			errResult = err
		}
	}
	return errResult
}
