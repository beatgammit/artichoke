package artichoke

import (
	"http"
)

func Static(root string) Middleware {
	return func (w http.ResponseWriter, r *http.Request, d Data) bool {
		http.ServeFile(w, r, root + r.RawURL)
		return true
	}
}
