package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Post struct {
	URL   string
	Date  string
	Desc  string
	Tags  []string
	Text  string
	Title string
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

func load_posts(path string, markdown goldmark.Markdown) (map[string][]Post, map[string]Post, []Post) {
	list_posts := make([]Post, 0)
	posts := make(map[string]Post)
	tags := make(map[string][]Post)
	files, _ := os.ReadDir(path)

	var post Post
	for i := range files {
		name := files[i].Name()
		data, _ := os.ReadFile(path + name)

		if !bytes.Equal(data[:4], []byte("---\n")) {
			continue
		}

		ix := bytes.Index(data[4:], []byte("---\n")) + 4
		post = Post{}
		yaml.Unmarshal(data[4:ix], &post)

		r, _ := regexp.Compile("#.*")

		var text bytes.Buffer
		markdown.Convert(data[ix+4:], &text)
		post.URL = "/" + path + strings.Split(name, ".")[0]
		post.Text = text.String()
		post.Title = r.FindString(string(data))[2:]

		posts[strings.Split(name, ".")[0]] = post

		for _, v := range post.Tags {
			tags[v] = append(tags[v], post)
		}

		list_posts = append(list_posts, post)
	}

	sort.SliceStable(list_posts, func(i, j int) bool {
		return list_posts[i].Date > list_posts[j].Date
	})

	return tags, posts, list_posts
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

	tags, posts, list_posts := load_posts("posts/", markdown)

	e := echo.New()

	e.Static("/static", "static")
	newTemplateRenderer(e, "templates/*.html")

	var bufPosts bytes.Buffer
	e.Renderer.Render(&bufPosts, "posts", map[string]any{"posts": list_posts}, nil)

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", map[string]any{"posts": list_posts[:4]})
	})

	e.GET("/unity-radionica", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "https://drive.google.com/drive/folders/1WzqJanDBeF2VAcEANw29mT-wFIzJFPlY?usp=sharing") // I know Google Drive shit
	})

	e.GET("/posts/:post", func(c echo.Context) error {
		post := c.Param("post")

		post_data, ok := posts[post]
		if ok {
			return c.Render(http.StatusOK, "post", post_data)
		}

		return c.Render(http.StatusOK, "404", nil)
	})

	e.GET("/posts", func(c echo.Context) error {
		/*
			res := map[string]any{
				"posts": list_posts,
			}
			return c.Render(http.StatusOK, "posts", res)
		*/
		return c.HTML(http.StatusOK, bufPosts.String())
	})

	e.GET("/tags/:tag", func(c echo.Context) error {
		tag := c.Param("tag")

		tag_posts, ok := tags[tag]
		if ok {
			return c.Render(http.StatusOK, "tags", map[string]any{"tag": tag, "posts": tag_posts})
		}

		return c.Render(http.StatusOK, "404", nil)
	})

	e.Logger.Fatal(e.Start(":8080"))

}
