package ingress

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/edge/ingress/local"
)

// LocalCommand provides context for the "local" command
type LocalCommand struct{}

func (d *LocalCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		panic(err)
	}

	if err := local.Add(mgr); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Local returns the "local" subcommand.
func Local(app *kingpin.CmdClause) {
	c := &LocalCommand{}
	app.Command("local", "Start the local Ingress operator").Action(c.run)
}
