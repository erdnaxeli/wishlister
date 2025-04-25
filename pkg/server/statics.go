package server

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed static
var statics embed.FS

func setStatics(e *echo.Echo) {
	subFS := echo.MustSubFS(statics, "static/css")
	e.StaticFS("/css", subFS)
}
