package utils

import (
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, code int, err error, message string) error {
	// set the error message into context
	c.Set("err", err.Error())
	// return the response
	return c.JSON(code, map[string]string{"message": message})
}
