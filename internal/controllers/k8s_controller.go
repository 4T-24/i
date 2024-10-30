package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/ctf"
	"instancer/internal/discord"
	"instancer/internal/env"
	"strconv"
	"sync"
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
	skip int

	challenges          map[string]client.Object
	ctfdChallengesSpecs map[string]*v1.ChallengeSpec

	CtfClient *ctf.Client

	chrono.TaskScheduler
	tasks map[string]chrono.ScheduledTask
}

var reconcilerMutex sync.Mutex

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;create;delete;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=create;delete;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;create
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;get;create;delete;watch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;delete
// +kubebuilder:rbac:groups=getambassador.io,resources=mappings,verbs=create;delete
// +kubebuilder:rbac:groups=getambassador.io,resources=tcpmappings,verbs=create;delete
// +kubebuilder:rbac:groups=i.4ts.fr,resources=challenges,verbs=get;list;watch
// +kubebuilder:rbac:groups=i.4ts.fr,resources=instancedchallenges,verbs=get;list;watch
// +kubebuilder:rbac:groups=i.4ts.fr,resources=oracleinstancedchallenges,verbs=get;list;watch
func (r *InstancierReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reconcilerMutex.Lock()
	defer reconcilerMutex.Unlock()

	_ = log.FromContext(ctx)

	if r.skip > 0 {
		r.skip--
		logrus.WithField("skips_left", r.skip).Info("Skipped init reconciliation, we already did it in the first occurence when fetching all.")
		return ctrl.Result{}, nil
	}

	// Init map if not init
	if !r.init {
		r.Reinit()
	}

	res, err := r.Register(ctx, req)
	if err != nil {
		return res, err
	}

	// Register the challenge onto CTFd (and ask for reconciliation if fail)
	return res, r.ReconcileCTFd()
}

func (r *InstancierReconciler) Register(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var challenge v1.Challenge
	err := r.Get(ctx, req.NamespacedName, &challenge)
	if err == nil {
		challenge.Status.Phase = "Syncing"
		r.Status().Update(context.Background(), &challenge)
		r.RegisterChallenge(&challenge)
		return ctrl.Result{}, nil
	}

	// Get the ressource and store it in our maps
	var instancedChallenge v1.InstancedChallenge
	err = r.Get(ctx, req.NamespacedName, &instancedChallenge)
	if err == nil {
		instancedChallenge.Status.Phase = "Syncing"
		r.Status().Update(context.Background(), &instancedChallenge)
		r.RegisterChallenge(&instancedChallenge)
		return ctrl.Result{}, nil
	}

	var oracleInstancedChallenge v1.OracleInstancedChallenge
	err = r.Get(ctx, req.NamespacedName, &oracleInstancedChallenge)
	if err == nil {
		oracleInstancedChallenge.Status.Phase = "Syncing"
		r.Status().Update(context.Background(), &oracleInstancedChallenge)
		r.RegisterChallenge(&oracleInstancedChallenge)
		return ctrl.Result{}, nil
	}

	var globallyInstancedChallenge v1.GloballyInstancedChallenge
	err = r.Get(ctx, req.NamespacedName, &globallyInstancedChallenge)
	if err == nil {
		globallyInstancedChallenge.Status.Phase = "Syncing"
		r.Status().Update(context.Background(), &globallyInstancedChallenge)
		r.RegisterChallenge(&globallyInstancedChallenge)
		return ctrl.Result{}, nil
	}

	logrus.WithField("name", req.Name).Warn("Ressource deleted, rechecking")
	r.Reinit()

	return ctrl.Result{}, nil
}

func (r *InstancierReconciler) ReconcileCTFd() error {
	logrus.WithField("challenges", len(r.ctfdChallengesSpecs)).Info("Reconciling CTFd with challenges")
	errors, err := r.CtfClient.ReconcileChallenge(r.ctfdChallengesSpecs)
	if err != nil {
		logrus.WithError(err).Error("Failed to reconcile challenges")
		return err
	}

	cfg := env.Get()

	for name, challObj := range r.challenges {
		var oldError, newError string

		err, found := errors[name]
		switch chall := challObj.(type) {
		case *v1.Challenge:
			var obj = v1.Challenge{
				ObjectMeta: chall.ObjectMeta,
			}
			r.Get(context.Background(), client.ObjectKeyFromObject(&obj), &obj)
			oldError = obj.Status.Error
			obj.Status.Error = ""
			obj.Status.Phase = "Synced"
			if found {
				obj.Status.Error = err.Error()
				obj.Status.Phase = "Error"
			}
			newError = obj.Status.Error
			r.Status().Update(context.Background(), &obj)
		case *v1.InstancedChallenge:
			var obj = v1.InstancedChallenge{
				ObjectMeta: chall.ObjectMeta,
			}
			r.Get(context.Background(), client.ObjectKeyFromObject(&obj), &obj)
			oldError = obj.Status.Error
			obj.Status.Error = ""
			obj.Status.Phase = "Synced"
			if found {
				obj.Status.Error = err.Error()
				obj.Status.Phase = "Error"
			}
			newError = obj.Status.Error
			r.Status().Update(context.Background(), &obj)
		case *v1.GloballyInstancedChallenge:
			var obj = v1.GloballyInstancedChallenge{
				ObjectMeta: chall.ObjectMeta,
			}
			r.Get(context.Background(), client.ObjectKeyFromObject(&obj), &obj)
			oldError = obj.Status.Error
			obj.Status.Error = ""
			obj.Status.Phase = "Synced"
			if found {
				obj.Status.Error = err.Error()
				obj.Status.Phase = "Error"
			}
			newError = obj.Status.Error
			r.Status().Update(context.Background(), &obj)
		case *v1.OracleInstancedChallenge:
			var obj = v1.OracleInstancedChallenge{
				ObjectMeta: chall.ObjectMeta,
			}
			r.Get(context.Background(), client.ObjectKeyFromObject(&obj), &obj)
			oldError = obj.Status.Error
			obj.Status.Error = ""
			obj.Status.Phase = "Synced"
			if found {
				obj.Status.Error = err.Error()
				obj.Status.Phase = "Error"
			}
			newError = obj.Status.Error
			r.Status().Update(context.Background(), &obj)
		}

		if cfg.Discord.WebhookToken != "" && oldError != newError {
			if newError != "" {
				discord.SendMessage(fmt.Sprintf("%s status changed to error", name), fmt.Sprintf("Got error when reconciling challenge : %v", newError), 16711680)
			} else {
				discord.SendMessage(fmt.Sprintf("%s status changed to resolved", name), fmt.Sprintf("Error was : %v", oldError), 65280)
			}
		}
	}

	return nil
}

// Reinit is the first run
func (r *InstancierReconciler) Reinit() {
	r.challenges = make(map[string]client.Object)
	r.ctfdChallengesSpecs = make(map[string]*v1.ChallengeSpec)
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
	err = r.List(context.Background(), &instancedChallenges)
	if err != nil {
		panic(err)
	}

	for _, challenge := range instancedChallenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var oracleInstancedChallenges v1.OracleInstancedChallengeList
	err = r.List(context.Background(), &oracleInstancedChallenges)
	if err != nil {
		panic(err)
	}

	for _, challenge := range oracleInstancedChallenges.Items {
		r.RegisterChallenge(&challenge)
	}

	var globallyInstancedChallenges v1.GloballyInstancedChallengeList
	err = r.List(context.Background(), &globallyInstancedChallenges)
	if err != nil {
		panic(err)
	}

	for _, challenge := range globallyInstancedChallenges.Items {
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
			timestamp := namespace.Labels["i.4ts.fr/stops-at-timestamp"]
			timestampI, err := strconv.Atoi(timestamp)
			if err != nil {
				logrus.WithError(err).WithField("namespace", namespace.Name).Warn("Failed to parse timestamp")
				continue
			}

			if time.Now().Unix() > int64(timestampI) {
				logrus.WithField("namespace", namespace.Name).Warn("Namespace expired")
				r.DeleteInstance(namespace.Labels["i.4ts.fr/challenge"], namespace.Labels["i.4ts.fr/instance"])
				continue
			}

			ttl := namespace.Labels["i.4ts.fr/ttl"]
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
		r.skip = len(r.challenges) - 1
	}
	r.init = true
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstancierReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&v1.Challenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	if err != nil {
		return err
	}
	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1.InstancedChallenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	if err != nil {
		return err
	}
	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1.GloballyInstancedChallenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	if err != nil {
		return err
	}
	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1.OracleInstancedChallenge{}, builder.WithPredicates(predicates())).
		Complete(r)
	if err != nil {
		return err
	}
	return nil
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
		UpdateFunc: func(e event.UpdateEvent) bool {
			var oldSpec, newSpec []byte

			switch newObject := e.ObjectNew.(type) {
			case *v1.Challenge:
				oldObject := e.ObjectOld.(*v1.Challenge)
				oldSpec, _ = json.Marshal(oldObject.Spec)
				newSpec, _ = json.Marshal(newObject.Spec)
			case *v1.InstancedChallenge:
				oldObject := e.ObjectOld.(*v1.InstancedChallenge)
				oldSpec, _ = json.Marshal(oldObject.Spec)
				newSpec, _ = json.Marshal(newObject.Spec)
			case *v1.GloballyInstancedChallenge:
				oldObject := e.ObjectOld.(*v1.GloballyInstancedChallenge)
				oldSpec, _ = json.Marshal(oldObject.Spec)
				newSpec, _ = json.Marshal(newObject.Spec)
			case *v1.OracleInstancedChallenge:
				oldObject := e.ObjectOld.(*v1.OracleInstancedChallenge)
				oldSpec, _ = json.Marshal(oldObject.Spec)
				newSpec, _ = json.Marshal(newObject.Spec)
			}
			return string(oldSpec) != string(newSpec)
		},
	}
}
