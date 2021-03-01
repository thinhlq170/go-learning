package main

import (
	"flag"
	"fmt"
	"os"

	b "example.com/main/cyoa/story"
)

// //Story type
// type Story map[string]Chapter

// //Chapter gens json structure
// type Chapter struct {
// 	Title   string   `json:"title"`
// 	Story   []string `json:"story"`
// 	Options []Option `json:"options"`
// }

// // Option is child type of AutoGenerated
// type Option struct {
// 	Text string `json:"text"`
// 	Arc  string `json:"arc"`
// }

// func jsonStory(fd io.Reader) (Story, error) {
// 	d := json.NewDecoder(fd)
// 	var story Story
// 	if err := d.Decode(&story); err != nil {
// 		return nil, err
// 	}
// 	return story, nil
// }

func main() {
	filename := flag.String("file", "gophers.json", "The JSON with the CYOA story")
	flag.Parse()

	fd, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	story, err := b.JsonStory(fd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", story)

}
