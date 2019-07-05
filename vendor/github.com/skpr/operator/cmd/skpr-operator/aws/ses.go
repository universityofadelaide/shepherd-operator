package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	controller "github.com/skpr/operator/internal/controller/aws/ses"
)

// SESCommand provides context for the "ses" command.
type SESCommand struct {
	params controller.Params
}

func (d *SESCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	var (
		iamClient = iam.New(session.New())
		sesClient = ses.New(session.New())
	)

	if err := controller.Add(mgr, iamClient, sesClient, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// SES returns the "ses" subcommand.
func SES(app *kingpin.CmdClause) {
	c := &SESCommand{}
	cmd := app.Command("ses", "Start the SES (SMTP) operator").Action(c.run)
	cmd.Flag("prefix", "Prefix to used for reference when creating AWS resources").Envar("SKPR_OPERATOR_SES_PREFIX").Required().StringVar(&c.params.Prefix)
	cmd.Flag("hostname", "Hostname of the AWS SES endpoint").Envar("SKPR_OPERATOR_SES_HOSTNAME").Default("email-smtp.us-west-2.amazonaws.com").StringVar(&c.params.Hostname)
	cmd.Flag("port", "Prefix to used for reference when creating AWS resources").Envar("SKPR_OPERATOR_SES_PORT").Default("1025").IntVar(&c.params.Port)
}
