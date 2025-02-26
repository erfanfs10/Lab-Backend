package middlewares

import (
	"bytes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CustomLogger() echo.MiddlewareFunc {

	customLoggerConfig := middleware.LoggerConfig{

		CustomTagFunc: func(c echo.Context, buf *bytes.Buffer) (int, error) {
			if errVal := c.Get("err"); errVal != nil {
				// Safely assert that the value is a string
				if errStr, ok := errVal.(string); ok {
					buf.WriteString(errStr)
				} else {
					buf.WriteString("error is not a string")
				}
			} else {
				buf.WriteString("")
			}
			return 0, nil
		},
		Format: `{"time":"${time_custom}",` + "\n" +
			`"remote_ip":"${remote_ip}",` + "\n" +
			`"host":"${host}",` + "\n" +
			`"user_agent":"${user_agent}",` + "\n" +
			`"method":"${method}",` + "\n" +
			`"uri":"${uri}",` + "\n" +
			`"error":"${custom}",` + "\n" +
			`"status":${status}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}
	return middleware.LoggerWithConfig(customLoggerConfig)
}
