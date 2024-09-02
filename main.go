package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jorgeav527/vehicle-model/configs"
	"github.com/jorgeav527/vehicle-model/routes"
)

// // Template struct for custom renderer
// type Template struct {
// 	tmpl *template.Template
// }

// // NewTemplate initializes a new Template instance with templates parsed from the views folder
// func NewTemplate() *Template {
// 	return &Template{
// 		tmpl: template.Must(template.ParseGlob("views/*.templ")),
// 	}
// }

// // Render method for Template struct
// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.tmpl.ExecuteTemplate(w, name, data)
// }

func main() {
	app := echo.New()
	app.Static("/public", "public")

	// Run the database
	configs.ConnectDB()

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	// Set the custom renderer
	// app.Renderer = NewTemplate()

	// Routes
	routes.HomeRoute(app)
	routes.VehicleModelRoute(app)

	// Start server
	fmt.Printf("Listening on %v\n", ":1323")
	app.Logger.Fatal(app.Start(":1323"))
}
