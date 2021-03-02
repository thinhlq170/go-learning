package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	b "example.com/main/cyoa/story"
)

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Stories}}
        <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
            <li>
                <a href="/{{.Arc}}">{{.Text}}</a>
            </li>
			{{end}}
        </ul>
    </body>
</html>
`

func main() {
	port := flag.Int("port", 3000, "the port to start CYOA web application on")
	filename := flag.String("file", "gophers.json", "The JSON with the CYOA story")
	flag.Parse()

	fd, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	//defer fd.Close()

	story, err := b.JSONStory(fd)
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New("").Parse(defaultHandlerTmpl))

	h := b.NewHandler(story,
		b.WithTemplate(tpl),
		b.WithPathFunc(pathFn),
	)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	mux.Handle("/", b.NewHandler(story))
	fmt.Printf("Starting the server at port: %d\n", *port)
	//fmt.Printf("%+v\n", story)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story"):]
}
