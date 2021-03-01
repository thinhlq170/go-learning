package story

import (
	"encoding/json"
	"io"
)

//Story type
type Story map[string]Chapter

//Chapter gens json structure
type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// Option is child type of AutoGenerated
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

//JsonStory parse json form from io reader
func JsonStory(fd io.Reader) (Story, error) {
	d := json.NewDecoder(fd)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
