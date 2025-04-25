package server

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed static
var statics embed.FS

func (s Server) setStatics() {
	subFS := echo.MustSubFS(statics, "static/css")
	s.e.StaticFS("/css", subFS)
}
