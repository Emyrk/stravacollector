package api

import "testing"

func TestRouteRequirementID(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   int64
		lite bool
		want int64
	}{
		{"full alternate", 784585, false, 891091},
		{"full original", 891091, false, 891091},
		{"unrelated", 3145259, false, 3145259},
		{"lite unchanged", 784585, true, 784585},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := routeRequirementID(tc.id, tc.lite); got != tc.want {
				t.Fatalf("got %d want %d", got, tc.want)
			}
		})
	}
}
