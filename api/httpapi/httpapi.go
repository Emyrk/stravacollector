package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"net/http"

	"github.com/Emyrk/strava/api/modelsdk"
)

func Write(_ context.Context, rw http.ResponseWriter, status int, response interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	// Pretty up JSON when testing.
	if flag.Lookup("test.v") != nil {
		enc.SetIndent("", "\t")
	}
	err := enc.Encode(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Read decodes JSON from the HTTP request into the value provided.
func Read(ctx context.Context, rw http.ResponseWriter, r *http.Request, value interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(value)
	if err != nil {
		Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Request body must be valid JSON.",
			Detail:  err.Error(),
		})
		return false
	}
	return true
}
