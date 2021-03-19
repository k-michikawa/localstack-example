package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	listenPort, listenPortOk := os.LookupEnv("LISTEN_PORT")
	if !listenPortOk {
		log.Panic("LISTEN_PORT is not found")
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	e.Logger.Fatal(e.Start(listenPort))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!!!!")
}
