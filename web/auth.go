package web

import (
	"instancer/internal/auth"
	"instancer/internal/controllers"

	"github.com/go-fuego/fuego"
)

func GetToken(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (string, error) {
	return func(c fuego.ContextNoBody) (string, error) {
		instanceId := c.PathParam("instanceId")

		return auth.Generate(instanceId)
	}
}
