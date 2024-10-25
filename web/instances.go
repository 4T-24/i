package web

import (
	"fmt"
	"instancer/internal/auth"
	"instancer/internal/controllers"
	eventwatcher "instancer/internal/eventWatcher"
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type APIResponse[K any] struct {
	Result string `json:"result"`
	Data   K      `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func IsInstanceSolved(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (bool, error) {
	return func(c fuego.ContextNoBody) (bool, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.IsInstanceSolved(challengeId, instanceId)
	}
}

func ListenInstance(reconciler *controllers.InstancierReconciler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		challengeId := r.PathValue("challengeId")
		instanceId := r.PathValue("instanceId")

		fmt.Println(challengeId)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not upgrade to websocket", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		_, p, err := conn.ReadMessage()
		fmt.Println(string(p))
		if err != nil {
			return
		}

		token := string(p)

		claims, err := auth.Verify(token)
		if err != nil {
			fmt.Println(err)
			return
		}

		if instanceId != claims.InstanceID {
			return
		}

		events, err := reconciler.GetEvents(challengeId, instanceId)
		if err != nil {
			return
		}

		worker := eventwatcher.NewWorker()
		defer close(worker.Channel)

		for _, eventType := range events {
			listener, found := eventwatcher.Listeners[eventType]
			if !found {
				listener = new(eventwatcher.Listener)
				eventwatcher.Listeners[eventType] = listener
			}
			listener.Add(worker)
			defer listener.Remove(worker)
		}

		status, err := reconciler.GetInstance(challengeId, instanceId)
		if err != nil {
			logrus.Error(err)
			return
		}
		conn.WriteJSON(status)

		for range worker.Channel {
			time.Sleep(1 * time.Second)
			status, _ := reconciler.GetInstance(challengeId, instanceId)
			conn.WriteJSON(status)
		}
	}
}

func GetInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.GetInstance(challengeId, instanceId)
	}
}

func GetGlobals(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) ([]*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) ([]*controllers.InstanceStatus, error) {
		return reconciler.GetGlobalInstances()
	}
}

func CreateInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.CreateInstance(challengeId, instanceId)
	}
}

func CreateGlobals(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) ([]*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) ([]*controllers.InstanceStatus, error) {
		return reconciler.CreateGlobalInstances()
	}
}

func DeleteInstance(reconciler *controllers.InstancierReconciler) func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
	return func(c fuego.ContextNoBody) (*controllers.InstanceStatus, error) {
		challengeId := c.PathParam("challengeId")
		instanceId := c.PathParam("instanceId")

		return reconciler.DeleteInstance(challengeId, instanceId)
	}
}
