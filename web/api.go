package web

import (
	"instancer/internal/controllers"

	"github.com/go-fuego/fuego"
)

func SetupServer(reconciler *controllers.InstancierReconciler) {
	s := fuego.NewServer()

	s.Addr = "0.0.0.0:8080"

	fuego.Get(s, "/api/v1/token/{instanceId}", GetToken(reconciler), FromCtfd)

	fuego.Get(s, "/api/v1/globals", GetGlobals(reconciler), FromAdmin)
	fuego.Post(s, "/api/v1/globals", CreateGlobals(reconciler), FromAdmin)

	fuego.Get(s, "/api/v1/{challengeId}/{instanceId}/is_solved", IsInstanceSolved(reconciler), FromCtfd)

	fuego.AllStd(s, "/api/v1/{challengeId}/{instanceId}/events", ListenInstance(reconciler))

	fuego.Get(s, "/api/v1/{challengeId}/{instanceId}", GetInstance(reconciler), FromCtfd)
	fuego.Post(s, "/api/v1/{challengeId}/{instanceId}", CreateInstance(reconciler), FromCtfd)
	fuego.Delete(s, "/api/v1/{challengeId}/{instanceId}", DeleteInstance(reconciler), FromCtfd)

	go s.Run()
}
