package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	cmdapp "github.com/skpr/operator/cmd/skpr-operator/app"
	cmdaws "github.com/skpr/operator/cmd/skpr-operator/aws"
	cmdbackup "github.com/skpr/operator/cmd/skpr-operator/backup"
	"github.com/skpr/operator/cmd/skpr-operator/edge"
	cmdmysql "github.com/skpr/operator/cmd/skpr-operator/mysql"
	cmdrestore "github.com/skpr/operator/cmd/skpr-operator/restore"
	cmdsearch "github.com/skpr/operator/cmd/skpr-operator/search"
	cmdversion "github.com/skpr/operator/cmd/skpr-operator/version"
)

func main() {
	logf.SetLogger(logf.ZapLogger(false))

	app := kingpin.New("skpr-operator", "Kubernetes native integrations for the Skipper platform")

	cmdversion.Command(app)

	grpAWS := app.Command("aws", "Operators for deploying AWS services eg. CloudFront/WAF/Certificates")
	cmdaws.CloudFront(grpAWS)
	cmdaws.CloudFrontInvalidation(grpAWS)
	cmdaws.Certificate(grpAWS)
	cmdaws.CertificateRequest(grpAWS)
	cmdaws.CertificateValidate(grpAWS)
	cmdaws.SES(grpAWS)

	edge.Operators(app)

	grpApp := app.Command("app", "Operators for deploying applications")
	cmdapp.Drupal(grpApp)

	grpMySQL := app.Command("mysql", "Operators for MySQL operations")
	cmdmysql.Shared(grpMySQL)
	cmdmysql.Image(grpMySQL)
	cmdmysql.ImageScheduled(grpMySQL)

	grpBackup := app.Command("backup", "Operators for Backup tasks")
	cmdbackup.Restic(grpBackup)
	cmdbackup.Scheduled(grpBackup)

	grpRestore := app.Command("restore", "Operators for Restore tasks")
	cmdrestore.Restic(grpRestore)

	grpSearch := app.Command("search", "Operators for search services")
	cmdsearch.Solr(grpSearch)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
