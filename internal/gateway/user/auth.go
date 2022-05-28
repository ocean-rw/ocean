package user

import (
	"net/http"
)

func (m *Mgr) Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// TODO auth
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
