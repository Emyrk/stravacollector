package river

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava"
	"github.com/Emyrk/strava/strava/stravalimit"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverlog"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog"
)

var (
	getActivitiesUnauthenticated = errors.New("get activities: unauthenticated athlete, please re-authenticate")
)

func (m *Manager) EnqueueForwardLoad(ctx context.Context, athleteID int64, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(ctx, ForwardLoadArgs{
		AthleteID: athleteID,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type ForwardLoadArgs struct {
	AthleteID int64 `db:"athlete_id" json:"athlete_id"`
}

func (ForwardLoadArgs) Kind() string { return "forward_load" }
func (a ForwardLoadArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Tags:     []string{fmt.Sprintf("%d", a.AthleteID)},
		Queue:    riverBackloadQueue,
		Priority: PriorityDefault,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	}
}

type ForwardLoadWorker struct {
	mgr *Manager
	river.WorkerDefaults[ForwardLoadArgs]
}

func (*ForwardLoadWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w *ForwardLoadWorker) Work(ctx context.Context, job *river.Job[ForwardLoadArgs]) error {
	logger := jobLogFields(w.mgr.logger, job).With().Int64("athlete_id", job.Args.AthleteID).Logger()
	now := time.Now()

	// Can the worker run? Are we at a strava rate limit?
	if err := w.stravaCheck(ctx, logger, now); err != nil {
		return err
	}

	// Get the athlete for oauth.
	athlogin, err := w.mgr.db.GetAthleteLogin(ctx, job.Args.AthleteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// A non-logged in athlete has nothing to load. The job has nothing to do
			_ = river.RecordOutput(ctx, "athlete has no authentication, skipping any loading")
			return nil
		}

		return fmt.Errorf("get athlete login: %w", err)
	}

	cli := strava.NewOAuthClient(w.mgr.oauthCfg.Client(ctx, athlogin.OAuthToken()))

	// Always fetch the latest load info.
	// TODO: Lock the row so another worker does not try to load the same athlete at the same time.
	athleteLoad, err := w.mgr.db.GetAthleteLoad(ctx, athlogin.AthleteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			athleteLoad, err = w.mgr.db.UpsertAthleteForwardLoad(ctx, database.UpsertAthleteForwardLoadParams{
				AthleteID: athlogin.AthleteID,
				// Choose some data in the far past
				ActivityTimeAfter: database.Timestamptz(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
				LastLoadComplete:  false,
				LastTouched:       database.Timestamptz(now),
				NextLoadNotBefore: database.Timestamptz(now),
			})
			if err != nil {
				return fmt.Errorf("upsert athlete forward load: %w", err)
			}
		} else {
			return fmt.Errorf("get athlete load: %w", err)
		}
	}

	initParams := strava.GetActivitiesParams{
		Page:    0,
		PerPage: 50,
		After:   athleteLoad.ActivityTimeAfter.Time,
	}

	activities, err := w.getActivities(ctx, cli, athlogin.AthleteID, initParams)
	if err != nil {
		if errors.Is(err, getActivitiesUnauthenticated) {
			return nil // User logged out, and stopped.
		}
		return err
	}

	// Insert rows
	err = w.mgr.db.InTx(func(store database.Store) error {
		params := database.UpsertAthleteForwardLoadParams{
			AthleteID:         athleteLoad.AthleteID,
			ActivityTimeAfter: database.Timestamptz(initParams.After),
			LastLoadComplete:  len(activities) == 0,
			LastTouched:       database.Timestamptz(now),
			// Just a little bump to let another user go next.
			NextLoadNotBefore: database.Timestamptz(now.Add(time.Millisecond * 200)),
		}

		if len(activities) == 0 {
			offset := rand.Intn(24 * 3)
			// Wait 7-10 days before trying again. Webhooks should capture all new activities.
			params.NextLoadNotBefore = database.Timestamptz(now.Add((time.Hour * 24 * 7) + (time.Hour * time.Duration(offset))))
		}

		for _, act := range activities {
			_, err := store.UpsertMapData(ctx, database.UpsertMapDataParams{
				ID:              act.Map.ID,
				SummaryPolyline: act.Map.SummaryPolyline,
			})
			if err != nil {
				return fmt.Errorf("upsert map summary (%d): %w", act.ID, err)
			}

			_, err = store.UpsertActivitySummary(ctx, database.UpsertActivitySummaryParams{
				ID:                 act.ID,
				AthleteID:          act.Athlete.ID,
				UploadID:           act.UploadID,
				ExternalID:         act.ExternalID,
				Name:               act.Name,
				Distance:           act.Distance,
				MovingTime:         act.MovingTime,
				ElapsedTime:        act.ElapsedTime,
				TotalElevationGain: act.TotalElevationGain,
				ActivityType:       act.Type,
				SportType:          act.SportType,
				WorkoutType:        act.WorkoutType,
				StartDate:          database.Timestamptz(act.StartDate),
				StartDateLocal:     database.Timestamptz(act.StartDateLocal),
				Timezone:           act.Timezone,
				UtcOffset:          act.UtcOffset,
				AchievementCount:   act.AchievementCount,
				KudosCount:         act.KudosCount,
				CommentCount:       act.CommentCount,
				AthleteCount:       act.AthleteCount,
				PhotoCount:         act.PhotoCount,
				MapID:              act.Map.ID,
				Trainer:            act.Trainer,
				Commute:            act.Commute,
				Manual:             act.Manual,
				Private:            act.Private,
				Flagged:            act.Flagged,
				GearID:             act.GearID,
				AverageSpeed:       act.AverageSpeed,
				MaxSpeed:           act.MaxSpeed,
				DeviceWatts:        act.DeviceWatts,
				HasHeartrate:       act.HasHeartrate,
				PrCount:            act.PrCount,
				TotalPhotoCount:    act.TotalPhotoCount,
				AverageHeartrate:   act.AverageHeartrate,
				MaxHeartrate:       act.MaxHeartrate,
			})
			if err != nil {
				return fmt.Errorf("upsert activity summary (%d): %w", act.ID, err)
			}

			// Backload bike rides for more deets
			if isBikeRide(act.Type) || isBikeRide(act.SportType) {
				_, err = w.mgr.EnqueueFetchActivity(ctx, FetchActivityArgs{
					Source:         database.ActivityDetailSourceBackload,
					ActivityID:     act.ID,
					AthleteID:      athleteLoad.AthleteID,
					HugelPotential: canBeHugel(act) || canBeHugelLite(act),
					OnHugelDates:   onHugelDate(act),
				}, activityJobPriority(act), func(j *river.InsertOpts) {
					// Delay by 5minutes.
					// We do this because sometimes strava loads 0 segments for a ride, and it takes some time
					// for segments to be populated. The ride might have just been uploaded.
					j.ScheduledAt = time.Now().Add(time.Minute * 5)
				})
				if err != nil {
					return fmt.Errorf("enqueue fetch activity: %w", err)
				}
			}

			if act.StartDate.After(params.ActivityTimeAfter.Time) {
				params.ActivityTimeAfter = database.Timestamptz(act.StartDate)
			}
		}

		riverlog.Logger(ctx).Info("Load step",
			slog.Time("before", initParams.After),
			slog.Time("time_after", params.ActivityTimeAfter.Time),
			slog.Bool("complete", params.LastLoadComplete),
			slog.Int("activities_loaded", len(activities)),
		)
		_, err := store.UpsertAthleteForwardLoad(ctx, params)
		if err != nil {
			return fmt.Errorf("update athlete load after loading: %w", err)
		}
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	if len(activities) > 0 {
		// Keep going until we have no more activities to load.
		return river.JobSnooze(time.Second * 5)
	}

	return nil
}

func (w *ForwardLoadWorker) stravaCheck(ctx context.Context, logger zerolog.Logger, now time.Time) error {
	iBuf, dBuf := int64(150), int64(500)

	if stravalimit.NextDailyReset(now) < time.Hour*3 {
		iBuf, dBuf = 80, 300
	}
	if stravalimit.NextDailyReset(now) < time.Hour*1 {
		iBuf, dBuf = 50, 150
	}
	if stravalimit.NextDailyReset(now) < time.Minute*20 {
		iBuf, dBuf = 50, 100
	}

	if ok, limitLogger := stravalimit.CanLogger(1, iBuf, dBuf, logger); !ok {
		w.mgr.rateLimitLogger.Do(func() {
			limitLogger.Error().
				Str("job", "forward_athlete_data").
				Msg("hitting strava rate limit, job will try again later")
		})

		return w.mgr.StravaSnooze(ctx)
	}
	return nil
}

func (w *ForwardLoadWorker) getActivities(ctx context.Context, cli *strava.Client, athlete int64, params strava.GetActivitiesParams) ([]strava.ActivitySummary, error) {
	activities, err := cli.GetActivities(ctx, params)
	if err != nil {
		if se := strava.IsAPIError(err); se != nil {
			if se.Response.StatusCode == http.StatusTooManyRequests {
				return nil, w.mgr.StravaSnooze(ctx)
			}

			if se.Response.StatusCode == 597 {
				return nil, w.mgr.StravaMaintaince(ctx, fmt.Sprintf("code=597"))
			}

			if se.Response.StatusCode == http.StatusUnauthorized || se.Response.StatusCode == http.StatusForbidden {
				// Delete unauthenticated athlete login
				_ = w.mgr.db.DeleteAthleteLogin(ctx, athlete)
				_ = river.RecordOutput(ctx, getActivitiesUnauthenticated.Error())
				return nil, getActivitiesUnauthenticated
			}
		}

		_ = river.RecordOutput(ctx, fmt.Sprintf("failed to fetch activities: %s", err.Error()))
		return nil, err
	}

	return activities, nil
}

func canBeHugel(summary strava.ActivitySummary) bool {
	return database.DistanceToMiles(summary.Distance) > 80 &&
		database.DistanceToFeet(summary.TotalElevationGain) > 8000
}

func canBeHugelLite(summary strava.ActivitySummary) bool {
	return database.DistanceToMiles(summary.Distance) > 35 &&
		database.DistanceToFeet(summary.TotalElevationGain) > 3500
}

func onHugelDate(summary strava.ActivitySummary) bool {
	if summary.StartDate.Year() == 2024 && summary.StartDate.Month() == time.November {
		if summary.StartDate.Day() > 7 && summary.StartDate.Day() < 12 {
			return true
		}

	}
	return false
}

// isBikeRide covers the weird stuff like "VirtualRide", "EBikeRide", "MountainBikeRide"
func isBikeRide(act string) bool {
	act = strings.ToLower(act)
	if strings.Contains(act, "bike") {
		return true
	}
	if strings.Contains(act, "ride") {
		return true
	}
	return false
}

func activityJobPriority(summary strava.ActivitySummary) int {
	priority := PriorityLow
	if time.Since(summary.StartDate) < (time.Hour * 24 * 7) {
		// Add priority for recent activities
		priority = PriorityDefault
	}

	if database.DistanceToMiles(summary.Distance) > 80 &&
		database.DistanceToFeet(summary.TotalElevationGain) > 7000 &&
		time.Since(summary.StartDate) < (time.Hour*24*14) {
		// Recent big rides should be synced
		priority = PriorityHigh
	}

	return priority
}
