package mysql

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/mysql/imagescheduled"
)

// ImageScheduledCommand provides context for the "image-scheduled" command.
type ImageScheduledCommand struct{}

func (i *ImageScheduledCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := imagescheduled.Add(mgr); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// ImageScheduled returns the "image-scheduled" subcommand.
func ImageScheduled(app *kingpin.CmdClause) {
	c := &ImageScheduledCommand{}
	app.Command("image-scheduled", "Start the Scheduled MySQL Image operator").Action(c.run)
}
