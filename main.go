package main

import (
	"log"

	"github.com/labstack/echo"
)

func main() {
	route := echo.New()

	route.GET("/", Home)

	if err := route.Start(":4000"); err != nil {
		log.Println("Not Running Server A...", err.Error())
	}

}

func Home(c echo.Context) error {
	return c.JSON(200, "Certo")
}
