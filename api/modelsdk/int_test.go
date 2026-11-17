package modelsdk_test

import (
	"encoding/json"
	"testing"

	"github.com/Emyrk/strava/api/modelsdk"
)

func TestStringIntJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   modelsdk.StringInt
		text string
	}{
		{"zero", 0, "0"},
		{"legacy ID", 2661162, "2661162"},
		{"above JavaScript safe integer", 9007199254740993, "9007199254740993"},
		{"maximum int64", 9223372036854775807, "9223372036854775807"},
		{"minimum int64", -9223372036854775808, "-9223372036854775808"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			quoted := `"` + tc.text + `"`
			data, err := json.Marshal(tc.id)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != quoted {
				t.Fatalf("Marshal = %s, want %s", data, quoted)
			}
			for _, input := range []string{tc.text, quoted} {
				var got modelsdk.StringInt
				if err := json.Unmarshal([]byte(input), &got); err != nil {
					t.Fatalf("Unmarshal(%s): %v", input, err)
				}
				if got != tc.id {
					t.Fatalf("Unmarshal(%s) = %d, want %d", input, got, tc.id)
				}
			}
		})
	}
}

func TestStringIntJSONInvalid(t *testing.T) {
	for _, input := range []string{
		`9223372036854775808`, `"9223372036854775808"`,
		`-9223372036854775809`, `"-9223372036854775809"`,
		`1.5`, `"1.5"`, `"not an ID"`, `""`, `true`, `{}`, `[]`,
	} {
		t.Run(input, func(t *testing.T) {
			var id modelsdk.StringInt
			if err := json.Unmarshal([]byte(input), &id); err == nil {
				t.Fatalf("Unmarshal(%s) succeeded, want an error", input)
			}
		})
	}
}

func TestAPIIdentifiersMarshalAsStrings(t *testing.T) {
	const id modelsdk.StringInt = 9007199254740993
	for _, tc := range []struct {
		name   string
		value  any
		fields []string
	}{
		{"activity", modelsdk.ActivitySummary{ActivityID: id, AthleteID: id, UploadID: id}, []string{"activity_id", "athlete_id", "upload_id"}},
		{"athlete", modelsdk.AthleteSummary{AthleteID: id}, []string{"athlete_id"}},
		{"load", modelsdk.AthleteLoad{AthleteID: 9007199254740993}, []string{"athlete_id"}},
		{"eddington", modelsdk.Eddington{AthleteID: 9007199254740993}, []string{"athlete_id"}},
		{"effort", modelsdk.SegmentEffort{ActivityID: id, EffortID: id, SegmentID: id}, []string{"activity_id", "effort_id", "segment_id"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			for _, field := range tc.fields {
				if got := string(fields[field]); got != `"9007199254740993"` {
					t.Errorf("%s = %s, want a quoted, exact ID", field, got)
				}
			}
		})
	}
}
