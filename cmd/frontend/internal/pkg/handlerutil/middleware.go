package handlerutil

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRFMiddleware is HTTP middleware that helps prevent cross-site request forgery. To make your
// forms compliant you will have to submit the token via the X-Csrf-Token header, which is made
// available in the client-side context.
func CSRFMiddleware(next http.Handler, secure bool) http.Handler {
	return csrf.Protect(
		[]byte("e953612ddddcdd5ec60d74e07d40218c"),
		csrf.CookieName("csrf_token"),
		csrf.Path("/"),
		csrf.Secure(secure),
	)(next)
}

// HTTPSRedirectLoadBalanced is the same middleware as HTTPSRedirect, but will
// take precendence of the X-Forwarded-Proto header if specified.
func HTTPSRedirectLoadBalanced(next http.Handler) http.Handler {
	return httpsRedirect(next, func(r *http.Request) string {
		if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			return proto
		}
		return r.URL.Scheme
	})
}

// HTTPSRedirect is an HTTP middleware that will redirect non-HTTPS requests
// to HTTPS. It uses JS for the redirect because this is the most reliable
// solution if reverse proxies are involved.
func HTTPSRedirect(next http.Handler) http.Handler {
	return httpsRedirect(next, func(r *http.Request) string {
		return r.URL.Scheme
	})
}

func httpsRedirect(next http.Handler, getScheme func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if getScheme(r) != "https" {
			w.Write([]byte(`<script>window.location.protocol = "https:";</script>`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
