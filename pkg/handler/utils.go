package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func buildSSE(event, context string) string {
	var result string
	if len(event) != 0 {
		result = result + "event: " + event + "\n"
	}
	if len(context) != 0 {
		result = result + "data: " + context + "\n"
	}
	result = result + "\n"
	return result
}
