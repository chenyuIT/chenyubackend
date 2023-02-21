package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Web() {
	facades.Route.Get("/", func(ctx http.Context) {
		ctx.Response().Json(200, http.Json{
			"Hello": "Goravel",
		})
	})

	// DB
	dbController := controllers.NewDBController()
	facades.Route.Get("/db", dbController.Index)

	// Websocket
	websocketController := controllers.NewWebsocketController()
	facades.Route.Get("/ws", websocketController.Server)

	// Validation
	validationController := controllers.NewValidationController()
	facades.Route.Post("/validation/json", validationController.Json)
	facades.Route.Post("/validation/request", validationController.Request)
	facades.Route.Post("/validation/form", validationController.Form)
}
