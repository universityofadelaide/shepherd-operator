package backup

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/skpr/operator/cmd/skpr-operator/manager"
	"github.com/skpr/operator/internal/controller/backup/restic"
)

// ResticCommand provides context for the "restic" command.
type ResticCommand struct {
	params restic.Params
}

func (d *ResticCommand) run(c *kingpin.ParseContext) error {
	mgr, err := manager.New()
	if err != nil {
		return errors.Wrap(err, "new manager failed")
	}

	if err := restic.Add(mgr, d.params); err != nil {
		return errors.Wrap(err, "add to manager failed")
	}

	return mgr.Start(signals.SetupSignalHandler())
}

// Restic returns the "restic" subcommand.
func Restic(app *kingpin.CmdClause) {
	c := &ResticCommand{}
	cmd := app.Command("restic", "Start the Restic Backup operator").Action(c.run)

	cmd.Flag("bucket", "Name of S3 bucket").
		Envar("SKPR_BACKUP_RESTIC_BUCKET").
		Required().
		StringVar(&c.params.Pod.Bucket)
	cmd.Flag("aws-key-id", "The AWS access key ID").
		Envar("AWS_ACCESS_KEY_ID").
		Required().
		StringVar(&c.params.Pod.KeyID)
	cmd.Flag("aws-secret-key", "The AWS secret key").
		Envar("AWS_SECRET_ACCESS_KEY").
		Required().
		StringVar(&c.params.Pod.AccessKey)
	cmd.Flag("cpu", "CPU which will be assigned to a Pod executing the backup").
		Envar("SKPR_BACKUP_RESTIC_CPU").
		Default("100m").
		StringVar(&c.params.Pod.CPU)
	cmd.Flag("memory", "Memory which will be assigned to a Pod executing the backup").
		Envar("SKPR_BACKUP_RESTIC_MEMORY").
		Default("512Mi").
		StringVar(&c.params.Pod.Memory)
	cmd.Flag("image", "Image which is used for executing the restic command").
		Envar("SKPR_BACKUP_RESTIC_IMAGE").
		Default("docker.io/restic/restic:0.9.5").
		StringVar(&c.params.Pod.ResticImage)
	cmd.Flag("mysql-image", "Image which is used for executing the MySQL command").
		Envar("SKPR_BACKUP_RESTIC_MYSQL_IMAGE").
		Default("docker.io/skpr/mtk-mysql:latest").
		StringVar(&c.params.Pod.MySQLImage)
	cmd.Flag("workingdir", "Directory where backup steps will be executed").
		Envar("SKPR_BACKUP_RESTIC_WORKINGDIR").
		Default("/home/skpr").
		StringVar(&c.params.Pod.WorkingDir)
}
