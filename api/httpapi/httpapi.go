package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"net/http"
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
