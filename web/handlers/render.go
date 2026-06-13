package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Render renders an HTML page template wrapped inside the main layout
func Render(c *gin.Context, statusCode int, pageTemplate string, data gin.H) {
	tmpl, err := template.New("layout").ParseFiles(
		"web/views/layout.html",
		"web/views/"+pageTemplate,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}
	err = tmpl.ExecuteTemplate(c.Writer, "layout.html", data)
	if err != nil {
		_ = c.Error(err)
	}
}

// RenderPartial renders an HTML partial template without the layout (useful for HTMX)
func RenderPartial(c *gin.Context, statusCode int, partialTemplate string, data gin.H) {
	tmpl, err := template.ParseFiles("web/views/" + partialTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		_ = c.Error(err)
	}
}
