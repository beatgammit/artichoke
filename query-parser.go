package artichoke

import (
	"net/http"
)

func QueryParser() Middleware {
	return func(w http.ResponseWriter, r *http.Request, m Data) bool {
		m["query"] = r.URL.Query()
		return false
	}
}
