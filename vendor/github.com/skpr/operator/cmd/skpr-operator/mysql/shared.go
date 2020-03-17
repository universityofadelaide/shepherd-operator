package mysql

import (
	"io/ioutil"

	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/mysql/shared"
)

// SharedCommand provides context for the "shared" command.
type SharedCommand struct {
	// Path to the CA file, loaded and added to the params.
	caFile string
	// Params passed into the controller.
	params shared.Params
}

func (d *SharedCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if d.caFile != "" {
		data, err := ioutil.ReadFile(d.caFile)
		if err != nil {
			return errors.Wrap(err, "failed to load CAfile")
		}

		d.params.Connection.CA = string(data)
	}

	if err := shared.Add(mgr, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Shared returns the "shared" subcommand.
func Shared(app *kingpin.CmdClause) {
	c := &SharedCommand{}
	cmd := app.Command("shared", "Start the MySQL Shared operator").Action(c.run)
	cmd.Flag("provisioner-name", "Name of this provisioner").
		Envar("SKPR_MYSQL_PROVISIONER_NAME").
		Default("shared").
		StringVar(&c.params.ProvisionerName)
	cmd.Flag("hostname", "Hostname which will be used to provision databases").
		Envar("SKPR_MYSQL_HOSTNAME").
		StringVar(&c.params.Connection.Hostname)
	cmd.Flag("port", "Port which will be used to provision databases").
		Envar("SKPR_MYSQL_PORT").
		IntVar(&c.params.Connection.Port)
	cmd.Flag("username", "Username which will be used to provision databases").
		Envar("SKPR_MYSQL_USERNAME").
		StringVar(&c.params.Connection.Username)
	cmd.Flag("password", "Password which will be used to provision databases").
		Envar("SKPR_MYSQL_PASSWORD").
		StringVar(&c.params.Connection.Password)
	cmd.Flag("ca", "Path to a CA file used when provisioning databases").
		Envar("SKPR_MYSQL_CA").
		StringVar(&c.caFile)
}
