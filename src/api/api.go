package api

import (
	"github.com/labstack/echo"
	"github.com/vcokltfre/volcan/src/api/modules/meta"
)

func Start(bind string) error {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Server", "Volcan")
			return next(c)
		}
	})

	// Meta routes
	e.GET("/status", meta.GetStatus)

	return e.Start(bind)
}
