package main

import (
	"htmx-test/util"
	"htmx-test/views"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/:name", func(c echo.Context) error {
		comp := views.Default()
		return util.Render(c, comp)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
