package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/aws/certificatevalidate"
)

// CertificateValidateCommand provides context for the "certificate-validate" command.
type CertificateValidateCommand struct {
	params certificatevalidate.Params
}

func (d *CertificateValidateCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	client := route53.New(session.New())

	if err := certificatevalidate.Add(mgr, client, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// CertificateValidate returns the "certificate-validate" subcommand.
func CertificateValidate(app *kingpin.CmdClause) {
	c := &CertificateValidateCommand{}
	cmd := app.Command("certificate-validate", "Start the AWS CertificateValidate operator").Action(c.run)
	cmd.Flag("zone", "Zone identifier for Route 53").
		Envar("SKPR_OPERATOR_CERTIFICATE_VALIDATE_ZONE").
		Required().
		StringVar(&c.params.Zone)
	cmd.Flag("domain", "Domain which this operator will add validation records for").
		Envar("SKPR_OPERATOR_CERTIFICATE_VALIDATE_DOMAIN").
		Required().
		StringVar(&c.params.Domain)
	cmd.Flag("ttl", "TTL which will be applied to DNS entries").
		Envar("SKPR_OPERATOR_CERTIFICATE_VALIDATE_TTL").
		Default("300").
		Int64Var(&c.params.TTL)
}
