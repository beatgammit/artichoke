package artichoke

import (
	"net/http"
)

func Static(root string) Middleware {
	return func(w http.ResponseWriter, r *http.Request, d Data) bool {
		http.ServeFile(w, r, root + r.URL.Path)
		return true
	}
}
