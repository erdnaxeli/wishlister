package main

import "github.com/labstack/echo/v4"

func setStatics(e *echo.Echo) {
	e.Static("/css", "static/css")
}
