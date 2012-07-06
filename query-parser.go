package artichoke

import (
	"net/http"
	"net/url"
)

func GetQuery(d Data) (url.Values) {
	if q, ok := d.Get("query"); ok {
		return q.(url.Values)
	}

	return nil
}

func QueryParser() Middleware {
	return func(w http.ResponseWriter, r *http.Request, m Data) bool {
		m.Set("query", r.URL.Query())
		return false
	}
}
