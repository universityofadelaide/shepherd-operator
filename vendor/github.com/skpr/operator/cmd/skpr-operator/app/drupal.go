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
	Nginx DrupalCommandExporter
	FPM   DrupalCommandExporter
}

// DrupalCommandExporter contains admin details for running application exporters.
type DrupalCommandExporter struct {
	Image  string
	Port   string
	CPU    string
	Memory string
}

func (ctx *DrupalCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "cannot setup manager")
	}

	nginxCPU, err := resource.ParseQuantity(ctx.Nginx.CPU)
	if err != nil {
		return errors.Wrap(err, "cannot parse Nginx exporter CPU")
	}

	nginxMemory, err := resource.ParseQuantity(ctx.Nginx.Memory)
	if err != nil {
		return errors.Wrap(err, "cannot parse Nginx exporter memory")
	}

	fpmCPU, err := resource.ParseQuantity(ctx.FPM.CPU)
	if err != nil {
		return errors.Wrap(err, "cannot parse FPM exporter CPU")
	}

	fpmMemory, err := resource.ParseQuantity(ctx.FPM.Memory)
	if err != nil {
		return errors.Wrap(err, "cannot parse FPM exporter memory")
	}

	exporters := drupal.Exporters{
		Nginx: drupal.Exporter{
			Image:  ctx.Nginx.Image,
			Port:   ctx.Nginx.Port,
			CPU:    nginxCPU,
			Memory: nginxMemory,
		},
		FPM: drupal.Exporter{
			Image:  ctx.FPM.Image,
			Port:   ctx.FPM.Port,
			CPU:    fpmCPU,
			Memory: fpmMemory,
		},
	}

	if err := drupal.Add(mgr, exporters); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Drupal returns the "drupal" subcommand.
func Drupal(app *kingpin.CmdClause) {
	c := &DrupalCommand{}

	cmd := app.Command("drupal", "Start the Drupal operator").Action(c.run)

	cmd.Flag("nginx-exporter-image", "Image to use for exporting Nginx metrics").Envar("SKPR_EXPORTER_NGINX_IMAGE").Default("docker.io/skpr/nginx-exporter:v0.2.0").StringVar(&c.Nginx.Image)
	cmd.Flag("nginx-exporter-port", "Port to receive requests").Envar("SKPR_EXPORTER_NGINX_PORT").Default("9113").StringVar(&c.Nginx.Port)
	cmd.Flag("nginx-exporter-cpu", "CPU allowance for the Nginx exporter process").Envar("SKPR_EXPORTER_NGINX_CPU").Default("50m").StringVar(&c.Nginx.CPU)
	cmd.Flag("nginx-exporter-memory", "Memory allowance for the Nginx exporter process").Envar("SKPR_EXPORTER_NGINX_MEMORY").Default("96Mi").StringVar(&c.Nginx.Memory)

	cmd.Flag("fpm-exporter-image", "Image to use for exporting Nginx metrics").Envar("SKPR_EXPORTER_FPM_IMAGE").Default("docker.io/skpr/fpm-exporter:v1.0.0").StringVar(&c.FPM.Image)
	cmd.Flag("fpm-exporter-port", "Port to receive requests").Envar("SKPR_EXPORTER_FPM_PORT").Default("9253").StringVar(&c.FPM.Port)
	cmd.Flag("fpm-exporter-cpu", "CPU allowance for the Nginx exporter process").Envar("SKPR_EXPORTER_FPM_CPU").Default("50m").StringVar(&c.FPM.CPU)
	cmd.Flag("fpm-exporter-memory", "Memory allowance for the Nginx exporter process").Envar("SKPR_EXPORTER_FPM_MEMORY").Default("96Mi").StringVar(&c.FPM.Memory)
}
