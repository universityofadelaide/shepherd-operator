package solr

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	promlog "github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	searchv1beta1 "github.com/skpr/operator/pkg/apis/search/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	k8ssync "github.com/skpr/operator/pkg/utils/k8s/sync"
)

const (
	// ControllerName is an identifier for the controller.
	ControllerName = "solr-controller"

	// Prefix applied to Kubernetes resources.
	Prefix = "solr"

	// EnvSolrHeap is the name of the env var used to define the heap size.
	EnvSolrHeap = "SOLR_HEAP"
	// EnvSolrJavaMem required to tune Solr resource consumption.
	EnvSolrJavaMem = "SOLR_JAVA_MEM"
	// EnvJVMOpts required to tune Solr resource cosumption.
	EnvJVMOpts = "JVM_OPTS"

	// EnvSolrCore is the name of the env var used to define the core name.
	EnvSolrCore = "SOLR_CORE"

	// MountName is used when connecting a PersistentVolumeClaim to a VolumeMount.
	MountName = "data"

	// LabelAppName for discovery.
	LabelAppName = "app_name"
	// LabelAppType for discovery.
	LabelAppType = "app_type"
	// LabelAppLayer for discovery.
	LabelAppLayer = "app_layer"

	// ContainerInit used for identifying the initialisation container.
	ContainerInit = "init"
	// ContainerSolr used for identifying the long running Solr container.
	ContainerSolr = "solr"
)

// Add creates a new Solr Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params ReconcileParams) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params ReconcileParams) reconcile.Reconciler {
	return &ReconcileSolr{
		params: params,

		Client:   mgr.GetClient(),
		recorder: mgr.GetRecorder(ControllerName),
		scheme:   mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("solr-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &searchv1beta1.Solr{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &searchv1beta1.Solr{},
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &searchv1beta1.Solr{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileSolr{}

// ReconcileSolr reconciles a Solr object
type ReconcileSolr struct {
	client.Client
	params   ReconcileParams
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// ReconcileParams contains parameters to pass to reconcile loop.
type ReconcileParams struct {
	// Init container used to initialize a Solr core.
	Init ReconcileParamsInit
	// Image used to run Solr, tags are provided by the CustomResourceDefinition.
	Image string
	// Port which Solr will respond to requests.
	Port int
	// StorageClass used to provision storage.
	StorageClass string
	// Path which storage will be mounted.
	StorageMount string
}

// ReconcileParamsInit for initializing a Solr core.
type ReconcileParamsInit struct {
	// Image used to enforce permissions.
	Image string
	// Tag for the image used to enforce permissions.
	Tag string
	// User to enforce for a data directory.
	User string
}

// Reconcile reads that state of the cluster for a Solr object and makes changes based on the state read
// and what is in the Solr.Spec
// Automatically generate RBAC rules to allow the Controller to read and write statefulesets
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=search.skpr.io,resources=solrs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=search.skpr.io,resources=solrs/status,verbs=get;update;patch
func (r *ReconcileSolr) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	solr := &searchv1beta1.Solr{}

	err := r.Get(context.TODO(), request.NamespacedName, solr)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	log.Info("Syncing objects")

	status, err := r.Sync(log, solr)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Syncing status")

	err = r.SyncStatus(log, solr, status)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// Sync resources required to run Solr.
func (r *ReconcileSolr) Sync(log promlog.Logger, solr *searchv1beta1.Solr) (searchv1beta1.SolrStatus, error) {
	status := searchv1beta1.SolrStatus{
		Core: solr.Spec.Core,
		Port: r.params.Port,
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MountName,
			Namespace: solr.ObjectMeta.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: solr.Spec.Resources.Storage,
				},
			},
		},
	}

	// Empty storageClass should default to the cluster default storageclass.
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class-1
	if r.params.StorageClass != "" {
		pvc.Spec.StorageClassName = &r.params.StorageClass
	}

	var (
		name               = r.name(solr)
		replicas           = int32(1)
		timeoutGracePeriod = int64(300)
	)

	// @todo, Merge these with the Drupal operator as a helper function.
	labels := map[string]string{
		LabelAppName:  solr.ObjectMeta.Name,
		LabelAppType:  "solr",
		LabelAppLayer: "search",
	}

	envs, err := buildEnvVars(solr)
	if err != nil {
		return status, err
	}

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: solr.ObjectMeta.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				pvc,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Labels:    labels,
					Namespace: solr.ObjectMeta.Namespace,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &timeoutGracePeriod,
					InitContainers: []corev1.Container{
						// Our Solr containers run as the user "solr".
						// This container will ensure that the permissions are set.
						// Otherwise Solr will fail to boot in the first instance.
						{
							Name:            ContainerInit,
							Image:           fmt.Sprintf("%s:%s", r.params.Init.Image, r.params.Init.Tag),
							ImagePullPolicy: corev1.PullAlways,
							Command: []string{
								"/bin/bash", "-c",
							},
							Args: []string{
								fmt.Sprintf("chown -R %s:%s %s", r.params.Init.User, r.params.Init.User, r.params.StorageMount),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      MountName,
									MountPath: r.params.StorageMount,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    solr.Spec.Resources.CPU.Request,
									corev1.ResourceMemory: solr.Spec.Resources.Memory,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    solr.Spec.Resources.CPU.Limit,
									corev1.ResourceMemory: solr.Spec.Resources.Memory,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  ContainerSolr,
							Image: fmt.Sprintf("%s:%s", r.params.Image, solr.Spec.Version),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(r.params.Port),
								},
							},
							Env: envs,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									// https://cwiki.apache.org/confluence/display/solr/Ping
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(r.params.Port),
									},
								},
								InitialDelaySeconds: 300,
								TimeoutSeconds:      10,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    solr.Spec.Resources.CPU.Request,
									corev1.ResourceMemory: solr.Spec.Resources.Memory,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    solr.Spec.Resources.CPU.Limit,
									corev1.ResourceMemory: solr.Spec.Resources.Memory,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      MountName,
									MountPath: r.params.StorageMount,
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, statefulset, k8ssync.StatefulSet(solr, statefulset.Spec, r.scheme))
	if err != nil {
		return status, err
	}
	log.Infof("Synced StatefulSet %s with status: %s", statefulset.ObjectMeta.Name, result)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Labels:    labels,
			Namespace: solr.ObjectMeta.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: int32(r.params.Port),
				},
			},
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, service, k8ssync.Service(solr, service.Spec, r.scheme))
	if err != nil {
		return status, err
	}
	log.Infof("Synced Service %s with status: %s", service.ObjectMeta.Name, result)

	status.Host = service.Spec.ClusterIP

	return status, nil
}

// SyncStatus of the resources which have been provisioned.
func (r *ReconcileSolr) SyncStatus(log promlog.Logger, solr *searchv1beta1.Solr, status searchv1beta1.SolrStatus) error {
	if diff := deep.Equal(solr.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		solr.Status = status

		return r.Status().Update(context.TODO(), solr)
	}

	return nil
}

// Helper function to generate a Kubernetes object name for Solr.
func (r *ReconcileSolr) name(solr *searchv1beta1.Solr) string {
	return fmt.Sprintf("%s-%s", Prefix, solr.ObjectMeta.Name)
}

// Helper function to generate environment variables for a Solr StatefulSet.
func buildEnvVars(solr *searchv1beta1.Solr) ([]corev1.EnvVar, error) {
	memory := fmt.Sprintf("%dm", solr.Spec.Resources.Memory.Value()/1024/1024)

	envs := []corev1.EnvVar{
		{
			Name:  EnvSolrCore,
			Value: solr.Spec.Core,
		},
		{
			Name:  EnvSolrHeap,
			Value: memory,
		},
		{
			Name:  EnvSolrJavaMem,
			Value: fmt.Sprintf("-Xms%s -Xmx%s", memory, memory),
		},
		{
			Name:  EnvJVMOpts,
			Value: fmt.Sprintf("-Xms%s -Xmx%s", memory, memory),
		},
	}

	return envs, nil
}
