package web

import (
	"instancer/internal/auth"
	"instancer/internal/controllers"
	event_watcher "instancer/internal/event_watcher"
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

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not upgrade to websocket", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		token := string(p)

		claims, err := auth.Verify(token)
		if err != nil {
			return
		}

		if instanceId != claims.InstanceID {
			return
		}

		events, err := reconciler.GetEvents(challengeId, instanceId)
		if err != nil {
			return
		}

		worker := event_watcher.NewWorker()
		defer close(worker.Channel)

		for _, eventType := range events {
			listener, found := event_watcher.Listeners[eventType]
			if !found {
				listener = new(event_watcher.Listener)
				event_watcher.Listeners[eventType] = listener
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

		// While the connection is alive, we'll send the status every 30 seconds, to be sure we're not stucking the user
		go func() {
			for {
				time.Sleep(30 * time.Second)
				status, _ := reconciler.GetInstance(challengeId, instanceId)
				err := conn.WriteJSON(status)
				if err != nil {
					return
				}
			}
		}()

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
