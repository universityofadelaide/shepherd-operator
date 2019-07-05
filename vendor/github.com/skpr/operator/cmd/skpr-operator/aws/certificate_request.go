package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/aws/certificaterequest"
)

// CertificateRequestCommand provides context for the "certificate-request" command.
type CertificateRequestCommand struct {
	params certificaterequest.Params
}

func (d *CertificateRequestCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	client := acm.New(session.New())

	if err := certificaterequest.Add(mgr, client, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// CertificateRequest returns the "certificate-request" subcommand.
func CertificateRequest(app *kingpin.CmdClause) {
	c := &CertificateRequestCommand{}
	cmd := app.Command("certificate-request", "Start the AWS CertificateRequest operator").Action(c.run)
	cmd.Flag("prefix", "Prefix to used for reference when creating new distributions").Envar("SKPR_OPERATOR_CERTIFICATE_REQUEST_PREFIX").Required().StringVar(&c.params.Prefix)
}
