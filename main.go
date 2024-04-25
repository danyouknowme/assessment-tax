package main

import (
	"net/http"

	"github.com/danyouknowme/assessment-tax/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.POST("/tax/calculations", api.CalculateTax)
	e.Logger.Fatal(e.Start(":1323"))
}
