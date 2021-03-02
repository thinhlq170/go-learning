package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	linkParser "example.com/main/htmllinkparser/parser"
)

func main() {

	file := flag.String("file", "ex2.html", "file need to be parsed")

	fd, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	listLinks, err := linkParser.LinkParser(fd)
	if err != nil {
		log.Fatal(err)
	}

	if len(listLinks) > 0 {
		for index, link := range listLinks {
			fmt.Printf("Link %d {Link: %s, Text: %s}\n", index, link.Href, link.Text)
		}
	}
	//fmt.Printf("%v", listLinks)

}
