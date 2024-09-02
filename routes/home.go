package routes

import (
	"github.com/jorgeav527/vehicle-model/controllers"
	"github.com/labstack/echo/v4"
)

func HomeRoute(e *echo.Echo) {
	e.GET("", controllers.Home)
}
