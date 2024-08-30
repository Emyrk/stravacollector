package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (m *Manager) refreshViews(ctx context.Context) {
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

		var hugelDone, superDone time.Duration
		var hugelErr, superErr error

		wg.Add(1)
		go func() {
			hugelErr = m.DB.RefreshHugelActivities(ctx)
			hugelDone = time.Since(start)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			superErr = m.DB.RefreshSuperHugelActivities(ctx)
			superDone = time.Since(start)
			wg.Done()
		}()

		wg.Wait()
		logger.Info().
			AnErr("super_err", superErr).
			AnErr("hugel_err", hugelErr).
			Str("super_duration", fmt.Sprintf("%.3fs", superDone.Seconds())).
			Str("hugel_duration", fmt.Sprintf("%.3fs", hugelDone.Seconds())).
			Msg("refresh views")
	}
}
