/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	"github.com/universityofadelaide/shepherd-operator/controllers/extension/backup"
	"github.com/universityofadelaide/shepherd-operator/controllers/extension/backupscheduled"
	"github.com/universityofadelaide/shepherd-operator/controllers/extension/restore"
	"github.com/universityofadelaide/shepherd-operator/controllers/extension/sync"
	syncbackup "github.com/universityofadelaide/shepherd-operator/controllers/extension/sync/backup"
	syncrestore "github.com/universityofadelaide/shepherd-operator/controllers/extension/sync/restore"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(extensionv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a47560cf.shepherd",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	osclient, err := osv1client.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to get OpenShift client")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to get Kubernetes client")
		os.Exit(1)
	}

	if err = (&backup.Reconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		Recorder:  mgr.GetEventRecorderFor(backup.ControllerName),
		ClientSet: clientset,
		Params: backup.Params{
			ResourceRequirements: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_BACKUP_CPU")),
					corev1.ResourceMemory: resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_BACKUP_MEMORY")),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_BACKUP_CPU")),
					corev1.ResourceMemory: resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_BACKUP_MEMORY")),
				},
			},
			WorkingDir: os.Getenv("SHEPHERD_OPERATOR_BACKUP_WORKING_DIR"),
			MySQL: backup.MySQL{
				Image: os.Getenv("SHEPHERD_OPERATOR_BACKUP_MYSQL_IMAGE"),
			},
			AWS: backup.AWS{
				Endpoint:       os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_ENDPOINT"),
				BucketName:     os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_BUCKET_NAME"),
				Image:          os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_IMAGE"),
				FieldKeyID:     os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_KEY_ID"),
				FieldAccessKey: os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_ACCESS_KEY"),
				Region:         os.Getenv("SHEPHERD_OPERATOR_BACKUP_AWS_REGION"),
			},
			FilterByLabelAndValue: backup.FilterByLabelAndValue{
				Key:   os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_KEY"),
				Value: os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_VALUE"),
			},
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", backup.ControllerName)
		os.Exit(1)
	}

	if err = (&restore.Reconciler{
		Client:    mgr.GetClient(),
		OpenShift: osclient,
		Scheme:    mgr.GetScheme(),
		Recorder:  mgr.GetEventRecorderFor(restore.ControllerName),
		Params: restore.Params{
			ResourceRequirements: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_RESTORE_CPU")),
					corev1.ResourceMemory: resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_RESTORE_MEMORY")),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_RESTORE_CPU")),
					corev1.ResourceMemory: resource.MustParse(os.Getenv("SHEPHERD_OPERATOR_RESTORE_MEMORY")),
				},
			},
			WorkingDir: os.Getenv("SHEPHERD_OPERATOR_RESTORE_WORKING_DIR"),
			MySQL: restore.MySQL{
				Image: os.Getenv("SHEPHERD_OPERATOR_RESTORE_MYSQL_IMAGE"),
			},
			AWS: restore.AWS{
				Endpoint:       os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_ENDPOINT"),
				BucketName:     os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_BUCKET_NAME"),
				Image:          os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_IMAGE"),
				FieldKeyID:     os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_KEY_ID"),
				FieldAccessKey: os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_ACCESS_KEY"),
				Region:         os.Getenv("SHEPHERD_OPERATOR_RESTORE_AWS_REGION"),
			},
			FilterByLabelAndValue: restore.FilterByLabelAndValue{
				Key:   os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_KEY"),
				Value: os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_VALUE"),
			},
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", restore.ControllerName)
		os.Exit(1)
	}

	if err = (&backupscheduled.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor(backupscheduled.ControllerName),
		Params: backupscheduled.Params{
			FilterByLabelAndValue: backupscheduled.FilterByLabelAndValue{
				Key:   os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_KEY"),
				Value: os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_VALUE"),
			},
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", backupscheduled.ControllerName)
		os.Exit(1)
	}

	if err := sync.SetupWithManager(mgr, sync.Params{
		Restore: syncrestore.Params{
			FilterByLabelAndValue: syncrestore.FilterByLabelAndValue{
				Key:   os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_KEY"),
				Value: os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_VALUE"),
			},
		},
		Backup: syncbackup.Params{
			FilterByLabelAndValue: syncbackup.FilterByLabelAndValue{
				Key:   os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_KEY"),
				Value: os.Getenv("SHEPHERD_OPERATOR_FILTER_LABEL_VALUE"),
			},
		},
	}, osclient); err != nil {
		setupLog.Error(err, "unable to create sync controller")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
