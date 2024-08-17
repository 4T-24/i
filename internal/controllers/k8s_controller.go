package controllers

import (
	"context"
	v1 "instancer/api/v1"
	"instancer/internal/ctf"
	"strconv"
	"time"

	core "k8s.io/api/core/v1"

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
	ctfdChallenges map[string]*v1.ChallengeSpec

	CtfClient *ctf.Client

	chrono.TaskScheduler
	tasks map[string]chrono.ScheduledTask
}

//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;create;delete;watch
//+kubebuilder:rbac:groups="",resources=services,verbs=create;delete;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;get;create;delete;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;delete
//+kubebuilder:rbac:groups=getambassador.io,resources=mappings,verbs=create;delete
//+kubebuilder:rbac:groups=getambassador.io,resources=tcpmappings,verbs=create;delete
//+kubebuilder:rbac:groups=i.4ts.fr,resources=challenges,verbs=get;list;watch
//+kubebuilder:rbac:groups=i.4ts.fr,resources=instancedchallenges,verbs=get;list;watch
//+kubebuilder:rbac:groups=i.4ts.fr,resources=oracleinstancedchallenges,verbs=get;list;watch

func (r *InstancierReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Init map if not init
	if !r.init {
		r.Reinit()
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

	logrus.WithField("name", req.Name).Warn("Ressource deleted, rechecking")
	r.Reinit()

	return ctrl.Result{}, nil
}

func (r *InstancierReconciler) ReconcileCTFd() {
	logrus.WithField("challenges", len(r.ctfdChallenges)).Info("Reconciling CTFd with challenges")
	err := r.CtfClient.ReconcileChallenge(r.ctfdChallenges)
	if err != nil {
		logrus.WithError(err).Error("Failed to reconcile CTFd")
		return
	}
}

// Reinit is the first run
func (r *InstancierReconciler) Reinit() {
	r.challenges = make(map[string]client.Object)
	r.ctfdChallenges = make(map[string]*v1.ChallengeSpec)
	r.tasks = make(map[string]chrono.ScheduledTask)

	r.TaskScheduler = chrono.NewDefaultTaskScheduler()

	var challenges v1.ChallengeList
	err := r.List(context.Background(), &challenges)
	if err != nil {
		panic(err)
	}

	for _, challenge := range challenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var instancedChallenges v1.InstancedChallengeList
	err = r.List(context.Background(), &instancedChallenges, client.InNamespace("default"))
	if err != nil {
		panic(err)
	}

	for _, challenge := range instancedChallenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var oracleInstancedChallenges v1.OracleInstancedChallengeList
	err = r.List(context.Background(), &oracleInstancedChallenges, client.InNamespace("default"))
	if err != nil {
		panic(err)
	}

	for _, challenge := range oracleInstancedChallenges.Items {
		r.RegisterChallenge(&challenge)
	}

	// Get all namespaces
	var namespaces core.NamespaceList
	err = r.List(context.Background(), &namespaces)
	if err != nil {
		panic(err)
	}

	for _, namespace := range namespaces.Items {
		if by, found := namespace.Labels["app.kubernetes.io/managed-by"]; found && by == "atsi" {
			ttl := namespace.Labels["i.4ts.fr/ttl"]
			if ttl == "" {
				logrus.WithField("namespace", namespace.Name).Warn("No TTL found")
				continue
			}

			ttlI, err := strconv.Atoi(ttl)
			if err != nil {
				logrus.WithError(err).WithField("namespace", namespace.Name).Warn("Failed to parse TTL")
				continue
			}

			task, err := r.Schedule(func(ctx context.Context) {
				r.DeleteInstance(namespace.Labels["i.4ts.fr/challenge"], namespace.Labels["i.4ts.fr/instance"])
			}, chrono.WithTime(namespace.CreationTimestamp.Time.Add(time.Duration(ttlI)*time.Second)))
			if err == nil {
				r.tasks[namespace.Name] = task
			}
		}
	}

	if !r.init {
		logrus.Info("Init successful")
	}
	r.init = true
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstancierReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctrl.NewControllerManagedBy(mgr).
		For(&v1.Challenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	ctrl.NewControllerManagedBy(mgr).
		For(&v1.InstancedChallenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.OracleInstancedChallenge{}, builder.WithPredicates(predicates())).
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
