package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	route := echo.New()

	route.GET("/", Home)
	route.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if err := route.Start(":4000"); err != nil {
		log.Println("Not Running Server A...", err.Error())
	}

}

func Home(c echo.Context) error {
	return c.JSON(200, "Certo")
}
