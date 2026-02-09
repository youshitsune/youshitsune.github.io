package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Post struct {
	Date  string
	Title string
	URL   string
	Text  string
}

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func newTemplate(templates *template.Template) *Template {
	return &Template{templates}
}

func newTemplateRenderer(e *echo.Echo, paths ...string) {
	tmpl := &template.Template{}
	for i := range paths {
		template.Must(tmpl.ParseGlob(paths[i]))
	}
	t := newTemplate(tmpl)
	e.Renderer = t
}

func check_posts(path string, markdown goldmark.Markdown, posts *map[string]Post, list_posts *[]Post) {
	files, _ := os.ReadDir(path)
	if len(files) != len(*posts) {
		*posts, *list_posts = load_posts(path, markdown)
	}
}

func load_posts(path string, markdown goldmark.Markdown) (map[string]Post, []Post) {
	list_posts := make([]Post, 0)
	posts := make(map[string]Post)
	files, _ := os.ReadDir(path)

	var post Post
	for i := range files {
		name := files[i].Name()
		data, _ := os.ReadFile(path + name)
		ctx := strings.Split(string(data), "\n")
		date := ctx[0]
		title := ctx[1]
		var text bytes.Buffer
		markdown.Convert([]byte(strings.Join(ctx[1:], "\n")), &text)
		post = Post{
			Date:  date,
			Title: title[2:],
			URL:   "/" + path + strings.Split(name, ".")[0],
			Text:  text.String(),
		}

		posts[strings.Split(name, ".")[0]] = post
		list_posts = append(list_posts, post)

	}

	sort.SliceStable(list_posts, func(i, j int) bool {
		return list_posts[i].Date > list_posts[j].Date
	})

	return posts, list_posts
}

func main() {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	posts, list_posts := load_posts("posts/", markdown)

	e := echo.New()

	e.Static("/static", "static")
	newTemplateRenderer(e, "templates/*.html")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/posts/:post", func(c echo.Context) error {
		post := c.Param("post")

		post_data, ok := posts[post]
		if ok {
			res := map[string]any{
				"Data": post_data.Text,
			}
			return c.Render(http.StatusOK, "post", res)
		}

		return c.Render(http.StatusOK, "404", nil)
	})

	e.GET("/posts", func(c echo.Context) error {
		check_posts("posts/", markdown, &posts, &list_posts)

		res := map[string]any{
			"posts": list_posts,
		}
		return c.Render(http.StatusOK, "posts", res)
	})

	e.Logger.Fatal(e.Start(":8080"))

}
