package database

import "testing"

func TestHugelActivitiesView(t *testing.T) {
	tests := []struct {
		name      string
		routeYear int
		lite      bool
		want      string
	}{
		{name: "2025 full", routeYear: 2025, want: "hugel_activities_2025"},
		{name: "2025 lite", routeYear: 2025, lite: true, want: "lite_hugel_activities_2025"},
		{name: "2026 full temporarily uses 2025", routeYear: 2026, want: "hugel_activities_2025"},
		{name: "2026 lite temporarily uses 2025", routeYear: 2026, lite: true, want: "lite_hugel_activities_2025"},
		{name: "current full view", want: "hugel_activities"},
		{name: "current lite view", lite: true, want: "lite_hugel_activities"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hugelActivitiesView(test.routeYear, test.lite); got != test.want {
				t.Fatalf("unexpected view: got %q, want %q", got, test.want)
			}
		})
	}
}
