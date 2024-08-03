package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

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
	e := echo.New()

	e.Static("/static", "static")
	newTemplateRenderer(e, "templates/*.html")
	e.GET("/", func(c echo.Context) error {
		res := map[string]interface{}{
			"Dark": true,
		}
		return c.Render(http.StatusOK, "index", res)
	})

	e.GET("/light", func(c echo.Context) error {
		res := map[string]interface{}{
			"Dark": false,
		}
		return c.Render(http.StatusOK, "index", res)
	})

	e.GET("/dark", func(c echo.Context) error {
		res := map[string]interface{}{
			"Dark": true,
		}
		return c.Render(http.StatusOK, "index", res)
	})

	e.GET("/works", func(c echo.Context) error {
		return c.Render(http.StatusOK, "works", nil)
	})

	e.GET("/posts", func(c echo.Context) error {
		return c.Render(http.StatusOK, "posts", nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
