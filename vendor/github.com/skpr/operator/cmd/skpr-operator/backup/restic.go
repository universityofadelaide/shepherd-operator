package mysql

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

	cmd.Flag("bucket", "Name of S3 bucket").Envar("SKPR_BACKUP_RESTIC_BUCKET").Required().StringVar(&c.params.Pod.Bucket)
	cmd.Flag("aws-key-id", "Name of S3 bucket").Envar("AWS_ACCESS_KEY_ID").Required().StringVar(&c.params.Pod.KeyID)
	cmd.Flag("aws-secret-key", "Name of S3 bucket").Envar("AWS_SECRET_ACCESS_KEY").Required().StringVar(&c.params.Pod.AccessKey)
	cmd.Flag("cpu", "CPU which will be assigned to a Pod executing the backup").Envar("SKPR_BACKUP_RESTIC_CPU").Default("100m").StringVar(&c.params.Pod.CPU)
	cmd.Flag("memory", "Memory which will be assigned to a Pod executing the backup").Envar("SKPR_BACKUP_RESTIC_MEMORY").Default("512Mi").StringVar(&c.params.Pod.Memory)
	cmd.Flag("image", "Image which is used for executing the restic command").Envar("SKPR_BACKUP_RESTIC_IMAGE").Default("docker.io/restic/restic:0.9.5").StringVar(&c.params.Pod.ResticImage)
	cmd.Flag("mysql-image", "Image which is used for executing the MySQL command").Envar("SKPR_BACKUP_RESTIC_MYSQL_IMAGE").Default("docker.io/library/mariadb:10").StringVar(&c.params.Pod.MySQLImage)
	cmd.Flag("workingdir", "Directory where backup steps will be executed").Envar("SKPR_BACKUP_RESTIC_WORKINGDIR").Default("/home/skpr").StringVar(&c.params.Pod.WorkingDir)
	cmd.Flag("tag", "Tag to apply to backups").Envar("SKPR_BACKUP_RESTIC_TAG").Default("system").StringsVar(&c.params.Pod.Tags)
	cmd.Flag("starting-deadline", "How to long wait before marking a CronJob as failed").Envar("SKPR_BACKUP_RESTIC_STARTING_DEADLINE").Default("600").Int64Var(&c.params.CronJob.StartingDeadline)
	cmd.Flag("active-deadline", "How to long wait before marking a CronJob as failed").Envar("SKPR_BACKUP_RESTIC_ACTIVE_DEADLINE").Default("3600").Int64Var(&c.params.CronJob.ActiveDeadline)
	cmd.Flag("backoff-limit", "How many times to fail before marking a CronJob as failed").Envar("SKPR_BACKUP_RESTIC_BACKOFF_LIMIT").Default("2").Int32Var(&c.params.CronJob.BackoffLimit)
	cmd.Flag("success-history", "How successful Jobs to keep").Envar("SKPR_BACKUP_RESTIC_SUCCESS_HISTORY").Default("10").Int32Var(&c.params.CronJob.SuccessHistory)
	cmd.Flag("failed-history", "How failed Jobs to keep").Envar("SKPR_BACKUP_RESTIC_FAILED_HISTORY").Default("10").Int32Var(&c.params.CronJob.FailedHistory)
}
