package httpmw

import (
	"net/http"
	"net/url"
	"strings"
)

func NoWWW() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.Host, "www.") {
				// Just send back to the non-www version.
				rd, err := url.Parse("https://" + strings.TrimPrefix(r.Host, "www."))
				if err != nil {
					http.Redirect(w, r, "https://"+strings.TrimPrefix(r.Host, "www."), http.StatusTemporaryRedirect)
					return
				}
				rd.Path = r.URL.Path
				rd.RawQuery = r.URL.RawQuery

				http.Redirect(w, r, rd.String(), http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
