package meta

import "github.com/labstack/echo"

func GetStatus(c echo.Context) error {
	return c.JSON(200, Status{
		Status: "ok",
	})
}
