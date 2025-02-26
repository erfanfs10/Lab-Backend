package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func SeparateLogs() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("---------------------------------------------------------")
			return next(c)
		}
	}
}
