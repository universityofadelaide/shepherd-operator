package ingress

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/edge/ingress/cloudfront"
)

// CloudFrontCommand provides context for the "ingress" command
type CloudFrontCommand struct {
	params cloudfront.Params
}

func (d *CloudFrontCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return err
	}

	if err := cloudfront.Add(mgr, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// CloudFront returns the "cloudfront" subcommand.
func CloudFront(app *kingpin.CmdClause) {
	c := &CloudFrontCommand{}
	cmd := app.Command("cloudfront", "Start the CloudFront Ingress operator").Action(c.run)
	cmd.Flag("origin-endpoint", "Endpoint of the origin to forward requests").
		Envar("SKPR_OPERATOR_EDGE_INGRESS_ORIGIN_ENDPOINT").
		Required().
		StringVar(&c.params.OriginEndpoint)
	cmd.Flag("origin-policy", "Policy which will be used when connecting to the origin").
		Envar("SKPR_OPERATOR_EDGE_INGRESS_ORIGIN_POLICY").
		Default("https-only").
		StringVar(&c.params.OriginPolicy)
	cmd.Flag("origin-timeout", "Policy which will be used when connecting to the origin").
		Envar("SKPR_OPERATOR_EDGE_INGRESS_ORIGIN_TIMEOUT").
		Default("60").
		Int64Var(&c.params.OriginTimeout)

	cmd.Flag("contour-timeout", "Timeout policy Contour will use for requests").
		Envar("SKPR_OPERATOR_EDGE_INGRESS_CONTOUR_TIMEOUT").
		Default("300s").
		StringVar(&c.params.ContourTimeout)
}
