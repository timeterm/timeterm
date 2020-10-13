package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Templates struct {
	templates *template.Template
}

func Load() (*Templates, error) {
	templates, err := template.ParseGlob("templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	return &Templates{templates}, nil
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type ErrorData struct {
	Message string
}

func RenderError(c echo.Context, status int, msg string) error {
	return c.Render(status, "error", &ErrorData{
		Message: msg,
	})
}
