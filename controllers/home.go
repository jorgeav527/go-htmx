package controllers

import (
	"github.com/jorgeav527/vehicle-model/helpers"
	"github.com/jorgeav527/vehicle-model/views/layout"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	// data := "hola"
	return helpers.Render(c, layout.Home())
}
