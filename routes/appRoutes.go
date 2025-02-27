package routes

import (
	"github.com/erfanfs10/Lab-Backend/handlers"
	"github.com/labstack/echo/v4"
)

func AppRoutes(g *echo.Group) {
	g.GET("", handlers.Home)
	g.POST("convert/", handlers.Convert)
	g.GET("ws/", handlers.WebsocketHandler)
}
