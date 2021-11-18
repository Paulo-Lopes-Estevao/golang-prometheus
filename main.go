package main

import (
	"log"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	route := echo.New()
	p := prometheus.NewPrometheus("echo", nil)

	route.GET("/", func(c echo.Context) error {
		return c.JSON(200, "Hello world!")
	})

	p.Use(route)

	if err := route.Start(":9000"); err != nil {
		log.Println("Not Running Server A...", err.Error())
	}

}
