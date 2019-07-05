package edge

import (
	"github.com/alecthomas/kingpin"

	"github.com/skpr/operator/cmd/skpr-operator/edge/ingress"
)

// Operators under the edge group.
func Operators(app *kingpin.Application) {
	group := app.Command("edge", "Operators for deploying edge services eg. CDN/WAF")
	operator := group.Command("ingress", "Operators for deploying edge services eg. CDN/WAF")
	ingress.CloudFront(operator)
	ingress.Local(operator)
}
