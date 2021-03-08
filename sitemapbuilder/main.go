package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"example.com/main/htmllinkparser/parser"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var page string
	fmt.Printf("Pleaser enter wanted page: ")
	fmt.Scan(&page)

	_, err := url.ParseRequestURI(page)
	check(err)

	response, err := http.Get(page)

	links, err := parser.LinkParser(response.Body)

	//fmt.Println(links)
	if len(links) > 0 {
		for index, link := range links {
			fmt.Printf("Link %d {Link: %s, Text: %s}\n", index, link.Href, link.Text)
		}
	}

}
