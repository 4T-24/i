package controllers

import (
	"context"
	v1 "instancer/api/v1"

	"codnect.io/chrono"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type InstancierReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	init bool

	availableChallenges       map[string]*v1.Challenge
	availableOracleChallenges map[string]*v1.OracleChallenge

	chrono.TaskScheduler
	tasks map[string]chrono.ScheduledTask
}

//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;create;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=create;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;create;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;delete
//+kubebuilder:rbac:groups=voyager.appscode.com,resources=ingress,verbs=create;delete
//+kubebuilder:rbac:groups=i.4ts.fr,resources=challenges,verbs=get;list;watch
//+kubebuilder:rbac:groups=i.4ts.fr,resources=oraclechallenges,verbs=get;list;watch

func (r *InstancierReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Init map if not init
	if !r.init {
		r.Init()
	}

	// Get the ressource and store it in our maps
	var challenge v1.Challenge
	err := r.Get(ctx, req.NamespacedName, &challenge)
	if err != nil {
		var oraclechallenge v1.OracleChallenge
		err := r.Get(ctx, req.NamespacedName, &oraclechallenge)
		if err != nil {
			return ctrl.Result{}, err
		}
		r.RegisterChallenge(&oraclechallenge)

		return ctrl.Result{}, err
	}

	r.RegisterChallenge(&challenge)

	return ctrl.Result{}, nil
}

// Init is the first run
func (r *InstancierReconciler) Init() {
	r.availableChallenges = make(map[string]*v1.Challenge)
	r.availableOracleChallenges = make(map[string]*v1.OracleChallenge)
	r.tasks = make(map[string]chrono.ScheduledTask)

	r.TaskScheduler = chrono.NewDefaultTaskScheduler()

	var challenges v1.ChallengeList
	err := r.List(context.Background(), &challenges, client.InNamespace("default"))
	if err != nil {
		panic(err)
	}

	for _, challenge := range challenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var oraclechallenges v1.OracleChallengeList
	err = r.List(context.Background(), &oraclechallenges, client.InNamespace("default"))
	if err != nil {
		panic(err)
	}

	for _, challenge := range oraclechallenges.Items {
		r.RegisterChallenge(&challenge)
	}

	r.init = true
	logrus.Info("Init successful")
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstancierReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Challenge{}, builder.WithPredicates(predicates())).
		Complete(r)
}

func predicates() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(createEvent event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return false
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return true
		},
	}
}
