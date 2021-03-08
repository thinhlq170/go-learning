package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"example.com/main/htmllinkparser/parser"
	"example.com/main/sitemapbuilder/builder"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var page string
	fmt.Printf("Please enter wanted page: ")
	fmt.Scan(&page)

	_, err := url.ParseRequestURI(page)
	check(err)

	response, err := http.Get(page)

	links, err := parser.LinkParser(response.Body)

	errRes := builder.LinkBuilder(page, links, os.Stdout)
	check(errRes)
}
