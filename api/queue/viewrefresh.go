package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (m *Manager) refreshViews(ctx context.Context) {
	last2023 := time.Time{}

	logger := m.Logger.With().Str("job", "refresh_views").Logger()
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("View refresh ended")
			return
		default:
		}

		time.Sleep(time.Minute * 5)

		wg := sync.WaitGroup{}
		start := time.Now()

		var hugelDone, hugel2023Done, superDone time.Duration
		var hugelErr, hugel2023Err, superErr, hugelLiteErr error

		wg.Add(1)
		go func() {
			hugelErr = m.DB.RefreshHugelActivities(ctx)

			hugelLiteErr = m.DB.RefreshHugelLiteActivities(ctx)
			hugelDone = time.Since(start)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			superErr = m.DB.RefreshSuperHugelActivities(ctx)
			superDone = time.Since(start)
			wg.Done()
		}()

		if time.Since(last2023) > time.Hour {
			wg.Add(1)
			go func() {
				hugel2023Err = m.DB.RefreshHugel2023Activities(ctx)
				hugel2023Done = time.Since(start)
				wg.Done()
			}()
			last2023 = time.Now()
		}

		wg.Wait()
		logger.Info().
			AnErr("super_err", superErr).
			AnErr("hugel_err", hugelErr).
			AnErr("hugel2023_err", hugel2023Err).
			AnErr("hugel_lite_err", hugelLiteErr).
			Str("super_duration", fmt.Sprintf("%.3fs", superDone.Seconds())).
			Str("hugel_duration", fmt.Sprintf("%.3fs", hugelDone.Seconds())).
			Str("hugel2023_duration", fmt.Sprintf("%.3fs", hugel2023Done.Seconds())).
			Msg("refresh views")
	}
}
