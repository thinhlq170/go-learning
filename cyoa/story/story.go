package story

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
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

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

//HandlerOption sets directly value into a handler pointer
type HandlerOption func(h *handler)

// WithTemplate is an option to provide a custom template to
// be used when rendering stories.
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFunc is an option to provide a custom function
// for processing the story chapter from the incoming request.
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// NewHandler will construct an http.Handler that will render
// the story provided.
// The default handler will use the full path (minus the / prefix)
// as the chapter name, defaulting to "intro" if the path is
// empty. The default template creates option links that follow
// this pattern.
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

//Story type
type Story map[string]Chapter

//Chapter gens json structure
type Chapter struct {
	Title   string   `json:"title"`
	Stories []string `json:"story"`
	Options []Option `json:"options"`
}

// Option is child type of Chapter
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

//JSONStory parse json form from io reader
func JSONStory(fd io.Reader) (Story, error) {
	d := json.NewDecoder(fd)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
