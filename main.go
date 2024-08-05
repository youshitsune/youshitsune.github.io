package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

const HTMLHEAD = `
    <head>
        <meta name="viewport" content="width=device-width,initial-scale=1.0,minimum-scale=1"><title>Youshitsune</title>
        <meta property="og:title" content="Youshitsune">
        <meta property="og:type" content="website">


        <meta property="og:image" content="/static/img/avatar.png">

        <meta property="og:url" content="https://youshitsune.top/">
        <meta property="og:description" content="The homepage of Astatine. This website is a demo the of the Hugo theme Astatine.">
        <meta name="Description" property="description" content="The homepage of Astatine. This website is a demo the of the Hugo theme Astatine.">
        <link rel="stylesheet" href="/static/css/style.css" />
        <link rel="me" href="https://tilde.town/@youshitsune">
    </head>
`

const BACK = `
	<b><a href="/" id="back">⇐ HOME</a></b>
`

type Post struct {
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

	e.GET("/works", func(c echo.Context) error {
		return c.Render(http.StatusOK, "works", nil)
	})

	e.GET("/posts/:post", func(c echo.Context) error {
		post := c.Param("post")
		data, err := os.ReadFile("posts/" + post + ".md")
		if err != nil {
			var buf bytes.Buffer
			data, _ = os.ReadFile("posts/404.md")
			markdown.Convert(data, &buf)
			return c.HTML(http.StatusOK, HTMLHEAD+"<body id='posts'>"+BACK+buf.String()+"</body>")
		}

		var buf bytes.Buffer
		markdown.Convert(data, &buf)

		return c.HTML(http.StatusOK, HTMLHEAD+"<body id='posts'>"+BACK+buf.String()+"</body>")
	})

	e.GET("/posts", func(c echo.Context) error {
		var posts []Post

		files, _ := os.ReadDir("posts/")

		for i := range files {
			name := files[i].Name()
			data, _ := os.ReadFile("posts/" + name)
			ctx := strings.Split(string(data), "\n")
			title := ctx[0]
			var desc bytes.Buffer
			markdown.Convert([]byte(ctx[1]), &desc)

			posts = append(posts, Post{
				Title:       title[2:],
				Description: desc.String(),
				URL:         "/posts/" + strings.Split(name, ".")[0],
			})
		}

		res := map[string]interface{}{
			"posts": posts,
		}
		return c.Render(http.StatusOK, "posts", res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
