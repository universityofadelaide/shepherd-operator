package backup

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"github.com/skpr/operator/internal/controller/backup/scheduled"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
)

// ScheduledCommand provides context for the "backup-scheduled" command.
type ScheduledCommand struct{}

func (d *ScheduledCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := scheduled.Add(mgr); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Scheduled returns the "backup-scheduled" subcommand.
func Scheduled(app *kingpin.CmdClause) {
	c := &ScheduledCommand{}
	app.Command("scheduled", "Start the Scheduled Backup operator").Action(c.run)
}
