package web

import (
	"instancer/internal/controllers"

	"github.com/go-fuego/fuego"
)

func SetupServer(reconciler *controllers.InstancierReconciler) {
	s := fuego.NewServer()

	s.Addr = "0.0.0.0:80"

	fuego.Get(s, "/api/v1/{challengeId}/{instanceId}/is_solved", IsInstanceSolved(reconciler))
	fuego.Get(s, "/api/v1/{challengeId}/{instanceId}", GetInstance(reconciler))
	fuego.Post(s, "/api/v1/{challengeId}/{instanceId}", CreateInstance(reconciler))
	fuego.Delete(s, "/api/v1/{challengeId}/{instanceId}", DeleteInstance(reconciler))

	go s.Run()
}
