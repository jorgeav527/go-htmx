package routes

import (
	"github.com/jorgeav527/vehicle-model/controllers"
	"github.com/labstack/echo/v4"
)

func VehicleModelRoute(e *echo.Echo) {
	e.POST("/vehicle-model", controllers.CreateVehicleModel)
	e.GET("/vehicle-model/:vehicleModelId", controllers.GetAVehicleModel)
	e.Any("/vehicle-model/edit/:vehicleModelId", controllers.EditAVehicleModel)
	e.DELETE("/vehicle-model/:vehicleModelId", controllers.DeleteAVehicleModel)
	e.GET("/vehicle-model", controllers.GetAllVehicleModel)
	e.POST("/vehicle-model/search", controllers.SearchVehicleModel)
}
