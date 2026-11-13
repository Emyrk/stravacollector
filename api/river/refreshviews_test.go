package river

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Emyrk/strava/database"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRefreshViewsIncludes2026AndLiteDependencies(t *testing.T) {
	for _, latest := range []bool{true, false} {
		name := "all"
		if latest {
			name = "latest"
		}
		t.Run(name, func(t *testing.T) {
			db := &refreshStore{}
			worker := &RefreshViewsWorker{mgr: &Manager{db: db, logger: zerolog.Nop()}}
			err := worker.Work(context.Background(), &river.Job[RefreshViewsArgs]{
				JobRow: &rivertype.JobRow{}, Args: RefreshViewsArgs{Latest: latest},
			})
			require.NoError(t, err)
			// Work joins its goroutines before returning; the slice is no longer written.
			expected := []string{"2026", "2025", "2025-lite"}
			if !latest {
				expected = append(expected, "2024", "2024-lite", "2023", "super")
			}
			require.ElementsMatch(t, expected, db.calls)
			index := func(value string) int {
				for i, call := range db.calls {
					if call == value {
						return i
					}
				}
				return -1
			}
			require.Less(t, index("2025"), index("2025-lite"))
			if !latest {
				require.Less(t, index("2024"), index("2024-lite"))
			}
		})
	}
}

type refreshStore struct {
	database.Store
	mu       sync.Mutex
	calls    []string
	failures map[string]error
}

func (s *refreshStore) record(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, name)
	return s.failures[name]
}
func (s *refreshStore) RefreshHugel2026Activities(context.Context) error { return s.record("2026") }
func (s *refreshStore) RefreshHugel2025Activities(context.Context) error { return s.record("2025") }
func (s *refreshStore) RefreshHugelLite2025Activities(context.Context) error {
	return s.record("2025-lite")
}
func (s *refreshStore) RefreshHugel2024Activities(context.Context) error { return s.record("2024") }
func (s *refreshStore) RefreshHugelLite2024Activities(context.Context) error {
	return s.record("2024-lite")
}
func (s *refreshStore) RefreshHugel2023Activities(context.Context) error  { return s.record("2023") }
func (s *refreshStore) RefreshSuperHugelActivities(context.Context) error { return s.record("super") }

func TestRefreshViewsReturnsErrors(t *testing.T) {
	for _, latest := range []bool{false, true} {
		views := []string{"2026", "2025", "2025-lite"}
		mode := "latest"
		if !latest {
			mode = "all"
			views = append(views, "2024", "2024-lite", "2023", "super")
		}
		for _, view := range views {
			t.Run(mode+"/"+view, func(t *testing.T) {
				failure := errors.New("refresh failed")
				db := &refreshStore{failures: map[string]error{view: failure}}
				worker := &RefreshViewsWorker{mgr: &Manager{db: db, logger: zerolog.Nop()}}
				err := worker.Work(context.Background(), &river.Job[RefreshViewsArgs]{JobRow: &rivertype.JobRow{}, Args: RefreshViewsArgs{Latest: latest}})
				require.ErrorIs(t, err, failure)
			})
		}
	}
}
