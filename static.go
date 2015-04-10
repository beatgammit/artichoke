package artichoke

import (
	"net/http"
	"os"
	"path"
)

func Static(root string) Middleware {
	return func(w http.ResponseWriter, r *http.Request, d Data) bool {
		fPath := path.Join(root, path.Clean(r.URL.Path))

		// if the path doesn't exist, continue down the stack
		f, err := os.Open(fPath)
		if err != nil {
			if os.IsNotExist(err) {
				return false
			}
		} else {
			f.Close()
		}

		http.ServeFile(w, r, fPath)
		return true
	}
}
