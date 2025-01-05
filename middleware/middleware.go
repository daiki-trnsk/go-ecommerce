package middleware

import (
	"net/http"
	token "github.com/daiki-trnsk/go-ecommerce/tokens"
	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ClientToken := c.Request().Header.Get("token")
		if ClientToken == "" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No Authorization Header Provided"})
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err})
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		return next(c)
	}
}