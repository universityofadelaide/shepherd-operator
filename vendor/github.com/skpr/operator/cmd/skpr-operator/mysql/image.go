package mysql

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/mysql/image"
	mysqlimage "github.com/skpr/operator/pkg/utils/mysql/image"
)

// ImageCommand provides context for the "image" command.
type ImageCommand struct {
	params mysqlimage.GenerateParams
}

func (i *ImageCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := image.Add(mgr, i.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Image returns the "Image" subcommand.
func Image(app *kingpin.CmdClause) {
	c := &ImageCommand{}
	cmd := app.Command("image", "Start the MySQL Image operator").Action(c.run)

	cmd.Flag("dump-image", "Image which will be used to dump the database and sanitize").
		Envar("SKPR_MYSQL_IMAGE_DUMP_IMAGE").
		Default("docker.io/skpr/mtk-dump:latest").
		StringVar(&c.params.Dump.Image)
	cmd.Flag("dump-cpu", "CPU which will be given to the dump process").
		Envar("SKPR_MYSQL_IMAGE_DUMP_CPU").
		Default("250m").
		StringVar(&c.params.Dump.CPU)
	cmd.Flag("dump-memory", "Memory which will be given to the dump process").
		Envar("SKPR_MYSQL_IMAGE_DUMP_MEMORY").
		Default("512Mi").
		StringVar(&c.params.Dump.Memory)

	cmd.Flag("build-image", "Image which will be used to build the container").
		Envar("SKPR_MYSQL_IMAGE_BUILD_IMAGE").
		Default("docker.io/skpr/mtk-build:latest").
		StringVar(&c.params.Build.Image)
	cmd.Flag("build-cpu", "CPU which will be given to the build process").
		Envar("SKPR_MYSQL_IMAGE_BUILD_CPU").
		Default("250m").
		StringVar(&c.params.Build.CPU)
	cmd.Flag("build-memory", "Memory which will be given to the build process").
		Envar("SKPR_MYSQL_IMAGE_BUILD_MEMORY").
		Default("512Mi").
		StringVar(&c.params.Build.Memory)

	cmd.Flag("docker-configmap", "ConfigMap which contains Docker configuration (~/.docker/config.json)").
		Envar("SKPR_MYSQL_IMAGE_DOCKER_CONFIGMAP").
		StringVar(&c.params.Docker.ConfigMap)

	cmd.Flag("aws-secret-name", "Secret which contains AWS credentials").
		Envar("SKPR_MYSQL_IMAGE_AWS_SECRET_NAME").
		StringVar(&c.params.AWS.Secret)
	cmd.Flag("aws-key-id", "Secret key which contains AWS credentials").
		Envar("SKPR_MYSQL_IMAGE_AWS_KEY_ID").
		StringVar(&c.params.AWS.KeyID)
	cmd.Flag("aws-access-key", "Secret key which contains AWS credentials").
		Envar("SKPR_MYSQL_IMAGE_AWS_ACCESS_KEY").
		StringVar(&c.params.AWS.AccessKey)
}
