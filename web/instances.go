package web

import (
	"instancer/internal/controllers"

	"github.com/go-fuego/fuego"
)

func IsInstanceSolved(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (bool, error) {
	return func(c fuego.ContextNoBody) (bool, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.IsInstanceSolved(challengeId, instanceId)
	}
}

func GetInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.GetInstance(challengeId, instanceId)
	}
}

func CreateInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.CreateInstance(challengeId, instanceId)
	}
}

func DeleteInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.DeleteInstance(challengeId, instanceId)
	}
}
