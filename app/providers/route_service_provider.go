package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"chenyucms/app/http"
	"chenyucms/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
	//Add HTTP middlewares
	kernel := http.Kernel{}
	facades.Route().GlobalMiddleware(kernel.Middleware()...)
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	//Add routes
	routes.Web()
}
