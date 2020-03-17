package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	controller "github.com/skpr/operator/internal/controller/aws/cloudfront"
)

// CloudFrontCommand provides context for the "cloudfront" command.
type CloudFrontCommand struct {
	params controller.Params
}

func (d *CloudFrontCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	client := cloudfront.New(session.New())

	if err := controller.Add(mgr, client, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// CloudFront returns the "cloudfront" subcommand.
func CloudFront(app *kingpin.CmdClause) {
	c := &CloudFrontCommand{}
	cmd := app.Command("cloudfront", "Start the CloudFront operator").Action(c.run)
	cmd.Flag("prefix", "Prefix to used for reference when creating new distributions").
		Envar("SKPR_OPERATOR_CLOUDFRONT_PREFIX").
		Required().
		StringVar(&c.params.Prefix)

	cmd.Flag("logging-bucket", "Bucket which logging data will be pushed to by CloudFront").
		Envar("SKPR_OPERATOR_CLOUDFRONT_LOGGING_BUCKET").
		Required().
		StringVar(&c.params.LoggingBucket)
}
