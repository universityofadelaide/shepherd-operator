package app

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/app/drupal"
)

// DrupalCommand provides context for "drupal" command
type DrupalCommand struct {
	Metrics Metrics
}

// Exporter which is used by Prometheus.
type Exporter struct {
	Image  string
	Port   string
	CPU    string
	Memory string
}

// Metrics used for autoscaling.
type Metrics struct {
	FPM Metric
}

// Metric used for autoscaling.
type Metric struct {
	Name     string
	Image    string
	CPU      string
	Memory   string
	Protocol string
	Port     string
	Path     string
}

func (ctx *DrupalCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "cannot setup manager")
	}

	metrics := drupal.Metrics{
		FPM: drupal.MetricsFPM{
			Name:     ctx.Metrics.FPM.Name,
			Image:    ctx.Metrics.FPM.Image,
			CPU:      resource.MustParse(ctx.Metrics.FPM.CPU),
			Memory:   resource.MustParse(ctx.Metrics.FPM.Memory),
			Protocol: ctx.Metrics.FPM.Protocol,
			Port:     ctx.Metrics.FPM.Port,
			Path:     ctx.Metrics.FPM.Path,
		},
	}

	if err := drupal.Add(mgr, metrics); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Drupal returns the "drupal" subcommand.
func Drupal(app *kingpin.CmdClause) {
	c := &DrupalCommand{}

	cmd := app.Command("drupal", "Start the Drupal operator").Action(c.run)

	cmd.Flag("fpm-metrics-name", "Name of the metric to use for autoscaling").
		Envar("SKPR_FPM_METRICS_NAME").
		Default("phpfpm_active_processes").
		StringVar(&c.Metrics.FPM.Name)
	cmd.Flag("fpm-metrics-image", "Image to use for exporting FPM metrics").
		Envar("SKPR_FPM_METRICS_IMAGE").
		Default("skpr/fpm-metrics-adapter:sidecar-v0.0.2").
		StringVar(&c.Metrics.FPM.Image)
	cmd.Flag("fpm-metrics-cpu", "CPU allowance for the FPM metrics process").
		Envar("SKPR_FPM_METRICS_CPU").
		Default("20m").
		StringVar(&c.Metrics.FPM.CPU)
	cmd.Flag("fpm-metrics-memory", "Memory allowance for the FPM metrics process").
		Envar("SKPR_FPM_METRICS_MEMORY").
		Default("32Mi").
		StringVar(&c.Metrics.FPM.Memory)
	cmd.Flag("fpm-metrics-protocol", "Protocol to receive requests").
		Envar("SKPR_FPM_METRICS_PROTOCOL").
		Default("http").
		StringVar(&c.Metrics.FPM.Port)
	cmd.Flag("fpm-metrics-port", "Port to receive requests").
		Envar("SKPR_FPM_METRICS_PORT").
		Default("80").
		StringVar(&c.Metrics.FPM.Port)
	cmd.Flag("fpm-metrics-path", "Path to receive requests").
		Envar("SKPR_FPM_METRICS_PATH").
		Default("/metrics").
		StringVar(&c.Metrics.FPM.Port)
}
