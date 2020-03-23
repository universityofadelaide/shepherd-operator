package scheduled

import (
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"

	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

// GetScheduledTime from an annotation map.
func GetScheduledTime(annotations map[string]string) (*time.Time, error) {
	timeRaw := annotations[shpmetav1.ScheduledAnnotation]
	if len(timeRaw) == 0 {
		return nil, nil
	}

	timeParsed, err := time.Parse(time.RFC3339, timeRaw)
	if err != nil {
		return nil, err
	}

	return &timeParsed, nil
}

// GetNextSchedule for a scheduled spec/status.
func GetNextSchedule(scheduledSpec shpmetav1.ScheduledSpec, scheduledStatus shpmetav1.ScheduledStatus, creationTime, now time.Time) (time.Time, time.Time, error) {
	schedule, err := cron.ParseStandard(scheduledSpec.CronTab)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "unable to parse .Spec.Schedule.CronTab")
	}

	var earliestTime time.Time

	if scheduledStatus.LastExecutedTime != nil {
		earliestTime = scheduledStatus.LastExecutedTime.Time
	} else {
		earliestTime = creationTime
	}

	if scheduledSpec.StartingDeadlineSeconds != nil {
		// controller is not going to schedule anything below this point
		schedulingDeadline := now.Add(-time.Second * time.Duration(*scheduledSpec.StartingDeadlineSeconds))

		if schedulingDeadline.After(earliestTime) {
			earliestTime = schedulingDeadline
		}
	}

	if earliestTime.After(now) {
		return time.Time{}, schedule.Next(now), nil
	}

	var (
		lastMissed time.Time
		starts     = 0
	)

	for t := schedule.Next(earliestTime); !t.After(now); t = schedule.Next(t) {
		lastMissed = t

		starts++

		if starts > 100 {
			// We can't get the most recent times so just return an empty slice
			return time.Time{}, time.Time{}, errors.New("too many missed start times")
		}
	}

	return lastMissed, schedule.Next(now), nil
}
