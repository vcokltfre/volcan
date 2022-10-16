package api

import (
	"crypto/subtle"
	"os"

	"github.com/labstack/echo"
	"github.com/vcokltfre/volcan/src/api/models"
)

func authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(401, models.Error{
				Error: "No token provided",
				Code:  models.ErrorInvalidToken,
			})
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(os.Getenv("API_TOKEN"))) != 1 {
			return c.JSON(401, models.Error{
				Error: "Invalid token",
				Code:  models.ErrorInvalidToken,
			})
		}

		return next(c)
	}
}
