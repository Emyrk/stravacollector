package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog"
)

func jobLogFields[T river.JobArgs](logger zerolog.Logger, j *river.Job[T]) zerolog.Logger {
	lastAttempt := ""
	if j.AttemptedAt != nil {
		lastAttempt = j.AttemptedAt.String()
	}
	return logger.With().
		Int64("job_id", j.ID).
		Str("queue", j.Queue).
		Strs("tags", j.Tags).
		Int("err_count", j.Attempt).
		Str("last_error", lastAttempt).
		Int16("priority", int16(j.Priority)).
		Logger()
}

type failedJob[A river.JobArgs] struct {
	Job   *river.Job[A]
	Args  A
	Error string
}
