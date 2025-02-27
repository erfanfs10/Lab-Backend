package main

import (
	"github.com/erfanfs10/Lab-Backend/middlewares"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.Use(middlewares.CustomLogger())
	e.Use(middlewares.SeparateLogs())

	e.Logger.Fatal(e.Start(":8000"))
}
