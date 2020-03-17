package search

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/search/solr"
)

const (
	defaultInitImage    = "docker.io/skpr/solr"
	defaultInitTag      = "init"
	defaultInitUser     = "solr"
	defaultImage        = "docker.io/skpr/solr"
	defaultPort         = "8983"
	defaultStorageClass = "standard"
	defaultStorageMount = "/opt/solr/server/solr/mycores"
)

// SolrCommand provides context for the "solr" command.
type SolrCommand struct {
	params solr.ReconcileParams
}

func (d *SolrCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := solr.Add(mgr, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Solr returns the "shared" subcommand.
func Solr(app *kingpin.CmdClause) {
	c := &SolrCommand{}
	cmd := app.Command("solr", "Start the Apache Solr operator").Action(c.run)

	cmd.Flag("init-image", "Container image for the init container").
		Envar("SKPR_SEARCH_SOLR_INIT_IMAGE").
		Default(defaultInitImage).
		StringVar(&c.params.Init.Image)
	cmd.Flag("init-tag", "Container image tag for the init container").
		Envar("SKPR_SEARCH_SOLR_INIT_TAG").
		Default(defaultInitTag).
		StringVar(&c.params.Init.Tag)
	cmd.Flag("init-user", "The linux user that the solr container runs as - used to set up file permissions").
		Envar("SKPR_SEARCH_SOLR_INIT_USER").
		Default(defaultInitUser).
		StringVar(&c.params.Init.User)

	cmd.Flag("image", "Container image for the Solr container").
		Envar("SKPR_SEARCH_SOLR_IMAGE").
		Default(defaultImage).
		StringVar(&c.params.Image)
	cmd.Flag("port", "Port which applications will interact with Solr").
		Envar("SKPR_SEARCH_SOLR_PORT").
		Default(defaultPort).
		IntVar(&c.params.Port)
	cmd.Flag("storage-class", "The storage class to use for solr data volume").
		Envar("SKPR_SEARCH_SOLR_STORAGE_CLASS").
		Default(defaultStorageClass).
		StringVar(&c.params.StorageClass)
	cmd.Flag("storage-mount", "Path to solr data directory").
		Envar("SKPR_SEARCH_SOLR_STORAGE_MOUNT").
		Default(defaultStorageMount).
		StringVar(&c.params.StorageMount)
}
