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
	Date        string
	Title       string
	Description string
	URL         string
}

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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

	e := echo.New()

	e.Static("/static", "static")
	newTemplateRenderer(e, "templates/*.html")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/posts/:post", func(c echo.Context) error {
		post := c.Param("post")
		data, err := os.ReadFile("posts/" + post + ".md")
		if err != nil {
			var buf bytes.Buffer
			data, _ = os.ReadFile("posts/404.md")
			markdown.Convert(data, &buf)
			res := map[string]interface{}{
				"data": buf.String(),
			}
			return c.Render(http.StatusOK, "post", res)
		}

		var buf bytes.Buffer
		tmp := strings.Split(string(data), "\n")
		markdown.Convert([]byte(strings.Join(tmp[1:], "\n")), &buf)
		res := map[string]interface{}{
			"Data": buf.String(),
		}

		return c.Render(http.StatusOK, "post", res)
	})

	e.GET("/posts", func(c echo.Context) error {
		var posts []Post

		files, _ := os.ReadDir("posts/")

		for i := range files {
			name := files[i].Name()
			data, _ := os.ReadFile("posts/" + name)
			ctx := strings.Split(string(data), "\n")
			date := ctx[0]
			title := ctx[1]
			var desc bytes.Buffer
			markdown.Convert([]byte(ctx[2]), &desc)

			posts = append(posts, Post{
				Date:        date,
				Title:       title[2:],
				Description: desc.String(),
				URL:         "/posts/" + strings.Split(name, ".")[0],
			})
		}

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date > posts[j].Date
		})

		res := map[string]interface{}{
			"posts": posts,
		}
		return c.Render(http.StatusOK, "posts", res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
