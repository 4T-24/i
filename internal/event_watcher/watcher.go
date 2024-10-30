package event_watcher

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Event struct {
	Kind      string
	Name      string
	Namespace string
	Phase     string
	Status    string
}

func (e *Event) Key() string {
	if e.Namespace != "" {
		return fmt.Sprintf("%s/%s/%s", e.Kind, e.Namespace, e.Name)
	}
	return fmt.Sprintf("%s/%s", e.Kind, e.Name)
}

// Fonction pour trouver une condition par son type
func getDeployCondition(conditions []appsv1.DeploymentCondition, conditionType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i] // Retourne un pointeur vers la condition trouvée
		}
	}
	return nil // Retourne nil si la condition n'est pas trouvée
}

func Watch() {
	config := ctrl.GetConfigOrDie()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Créer des informers pour Pods et Namespaces
	factory := informers.NewSharedInformerFactory(clientset, time.Minute*10)

	deployInformer := factory.Apps().V1().Deployments().Informer()
	namespaceInformer := factory.Core().V1().Namespaces().Informer()

	// Configurer des handlers pour les événements
	deployInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)

			evt := &Event{
				Kind:      "deploy",
				Name:      deploy.Name,
				Namespace: deploy.Namespace,
			}
			l, found := Listeners[evt.Key()]
			if found {
				l.Send(evt)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deploy := newObj.(*appsv1.Deployment)
			cond := getDeployCondition(deploy.Status.Conditions, appsv1.DeploymentAvailable)
			if cond != nil {
				evt := &Event{
					Kind:      "deploy",
					Name:      deploy.Name,
					Namespace: deploy.Namespace,
					Status:    string(cond.Status),
				}
				l, found := Listeners[evt.Key()]
				if found {
					l.Send(evt)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)

			evt := &Event{
				Kind:      "deploy",
				Name:      deploy.Name,
				Namespace: deploy.Namespace,
			}
			l, found := Listeners[evt.Key()]
			if found {
				l.Send(evt)
			}
		},
	})

	namespaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			namespace := obj.(*corev1.Namespace)

			evt := &Event{
				Kind:  "namespace",
				Name:  namespace.Name,
				Phase: string(namespace.Status.Phase),
			}
			l, found := Listeners[evt.Key()]
			if found {
				l.Send(evt)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newNs := newObj.(*corev1.Namespace)

			evt := &Event{
				Kind:  "namespace",
				Name:  newNs.Name,
				Phase: string(newNs.Status.Phase),
			}
			l, found := Listeners[evt.Key()]
			if found {
				l.Send(evt)
			}
		},
		DeleteFunc: func(obj interface{}) {
			namespace := obj.(*corev1.Namespace)

			evt := &Event{
				Kind:  "namespace",
				Name:  namespace.Name,
				Phase: string(namespace.Status.Phase),
			}
			l, found := Listeners[evt.Key()]
			if found {
				l.Send(evt)
			}
		},
	})

	// Démarrer les informers
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// Bloquer l'exécution pour continuer à surveiller les événements
	<-stopCh
}
