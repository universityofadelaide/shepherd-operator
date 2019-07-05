package mysql

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/mysql/shared"
)

// SharedCommand provides context for the "shared" command.
type SharedCommand struct {
	name string
	conn shared.Connection
}

func (d *SharedCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := shared.Add(mgr, d.name, d.conn); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Shared returns the "shared" subcommand.
func Shared(app *kingpin.CmdClause) {
	c := &SharedCommand{}
	cmd := app.Command("shared", "Start the MySQL Shared operator").Action(c.run)
	cmd.Flag("provisioner-name", "Name of this provisioner").Envar("SKPR_MYSQL_PROVISIONER_NAME").Default("shared").StringVar(&c.name)
	cmd.Flag("hostname", "Hostname which will be used to provision databases").Envar("SKPR_MYSQL_HOSTNAME").StringVar(&c.conn.Hostname)
	cmd.Flag("port", "Port which will be used to provision databases").Envar("SKPR_MYSQL_PORT").IntVar(&c.conn.Port)
	cmd.Flag("username", "Username which will be used to provision databases").Envar("SKPR_MYSQL_USERNAME").StringVar(&c.conn.Username)
	cmd.Flag("password", "Password which will be used to provision databases").Envar("SKPR_MYSQL_PASSWORD").StringVar(&c.conn.Password)
}
