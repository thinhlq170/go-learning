package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"example.com/main/sitemapbuilder/builder"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	var page string
	fmt.Printf("Please enter wanted page: ")
	fmt.Scan(&page)

	_, err := url.ParseRequestURI(page)
	check(err)

	urlFlag := flag.String("url", page, "the url that you want to build a sitemap for")
	maxDepthFlag := flag.Int("maxDepth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	//response, err := http.Get(page)

	//links, err := parser.LinkParser(response.Body)

	// errRes := builder.LinkBuilder(page, links, os.Stdout)
	// check(errRes)

	type urlBuilder struct {
		XMLName xml.Name `xml:"url"`
		LOC     string   `xml:"loc"`
	}

	toXML := urlset{
		Xmlns: xmlns,
	}

	links := builder.Mybfs(*urlFlag, *maxDepthFlag)
	//linksBuilder := make([]urlBuilder, 0, len(links))
	for _, link := range links {
		//linksBuilder = append(linksBuilder, urlBuilder{LOC: link})
		toXML.Urls = append(toXML.Urls, loc{Value: link})
	}

	fmt.Printf(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
}
