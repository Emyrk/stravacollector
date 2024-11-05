package httpmw

import (
	"net/http"
	"strings"
)

func NoWWW() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.Host, "www.") {
				// Just send back to the non-www version.
				http.Redirect(w, r, "https://"+strings.TrimPrefix(r.Host, "www."), http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
