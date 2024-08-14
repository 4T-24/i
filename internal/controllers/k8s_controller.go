package controllers

import (
	"context"
	v1 "instancer/api/v1"
	"instancer/internal/ctf"

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

	challenges     map[string]client.Object
	ctfdChallenges []*v1.ChallengeSpec

	CtfClient *ctf.Client

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

	// Register the challenge onto CTFd
	defer r.ReconcileCTFd()

	var challenge v1.Challenge
	err := r.Get(ctx, req.NamespacedName, &challenge)
	if err == nil {
		if challenge.DeletionTimestamp != nil {
			r.UnregisterChallenge(&challenge)
			return ctrl.Result{}, nil
		}
		r.RegisterChallenge(&challenge)
		return ctrl.Result{}, nil
	}

	// Get the ressource and store it in our maps
	var instancedChallenge v1.InstancedChallenge
	err = r.Get(ctx, req.NamespacedName, &instancedChallenge)
	if err == nil {
		if instancedChallenge.DeletionTimestamp != nil {
			r.UnregisterChallenge(&instancedChallenge)
			return ctrl.Result{}, nil
		}
		r.RegisterChallenge(&instancedChallenge)
		return ctrl.Result{}, nil
	}

	var oracleInstancedChallenge v1.OracleInstancedChallenge
	err = r.Get(ctx, req.NamespacedName, &oracleInstancedChallenge)
	if err == nil {
		if oracleInstancedChallenge.DeletionTimestamp != nil {
			r.UnregisterChallenge(&oracleInstancedChallenge)
			return ctrl.Result{}, nil
		}
		r.RegisterChallenge(&oracleInstancedChallenge)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, err
}

func (r *InstancierReconciler) ReconcileCTFd() {
	err := r.CtfClient.ReconcileChallenge(r.ctfdChallenges)
	if err != nil {
		logrus.WithError(err).Error("Failed to reconcile CTFd")
		return
	}
}

// Init is the first run
func (r *InstancierReconciler) Init() {
	r.challenges = make(map[string]client.Object)
	r.tasks = make(map[string]chrono.ScheduledTask)

	r.TaskScheduler = chrono.NewDefaultTaskScheduler()

	var challenges v1.InstancedChallengeList
	err := r.List(context.Background(), &challenges, client.InNamespace("default"))
	if err != nil {
		panic(err)
	}

	for _, challenge := range challenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var oraclechallenges v1.OracleInstancedChallengeList
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
		Owns(&v1.InstancedChallenge{}, builder.WithPredicates(predicates())).
		Owns(&v1.OracleInstancedChallenge{}, builder.WithPredicates(predicates())).
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
