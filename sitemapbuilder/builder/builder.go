package builder

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
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

//Mybfs uses bfs algorithm to search all links from the page "urlStr" and returns slice of string
func Mybfs(urlStr string, maxDepth int) []string {

	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}

	for i := 0; i < maxDepth+1; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			links := getURL(url)
			for _, link := range links {
				nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for link := range seen {
		ret = append(ret, link)
	}
	return ret
}

func getURL(urlStr string) []string {
	res, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	reqURL := res.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	return filter(hrefs(res.Body, base), withPrefixFn(base))
}

func hrefs(r io.Reader, base string) []string {
	var ret []string

	links, err := parser.LinkParser(r)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			ret = append(ret, base+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			ret = append(ret, link.Href)
		}
	}

	return ret
}

func filter(links []string, checkPrefixFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if checkPrefixFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefixFn(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
