package aws

import (
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/aws/cloudfrontinvalidation"
)

// CloudFrontInvalidationCommand provides context for the "cloudfront-invalidation" command.
type CloudFrontInvalidationCommand struct{}

func (d *CloudFrontInvalidationCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	client := cloudfront.New(session.New())

	if err := cloudfrontinvalidation.Add(mgr, client); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// CloudFrontInvalidation returns the "cloudfront-invalidation" subcommand.
func CloudFrontInvalidation(app *kingpin.CmdClause) {
	c := &CloudFrontInvalidationCommand{}
	app.Command("cloudfront-invalidation", "Start the AWS CloudFrontInvalidation operator").Action(c.run)
}
