package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/aws/certificate"
)

// CertificateCommand provides context for the "certificate" command.
type CertificateCommand struct {
	params certificate.Params
}

func (d *CertificateCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := certificate.Add(mgr, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Certificate returns the "certificate" subcommand.
func Certificate(app *kingpin.CmdClause) {
	c := &CertificateCommand{}
	cmd := app.Command("certificate", "Start the AWS Certificate operator").Action(c.run)
	cmd.Flag("retention", "How many CertificateRequests to keep before cleaning up").
		Envar("SKPR_OPERATOR_CERTIFICATE_RETENTION").
		Default("5").
		IntVar(&c.params.Retention)
}
