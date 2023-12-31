package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("auth middleware")
		token := c.Request().Header.Get("Authorization")
		// TODO check if token is actually valid
		if token == "" {
			return c.Redirect(302, "/login")
		}
		fmt.Println(token)
		splitToken := strings.Split(token, "Bearer ")
		parsedToken := splitToken[1]
		c.Set("token", parsedToken)
		return next(c)
	}
}
