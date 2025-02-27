package main

import (
	"github.com/erfanfs10/Lab-Backend/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middlewares.CustomLogger())
	e.Use(middlewares.SeparateLogs())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST,
			echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin,
			echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Logger.Fatal(e.Start(":8000"))
}
