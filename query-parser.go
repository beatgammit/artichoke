package artichoke

import (
	"net/http"
	"net/url"
)

func GetQuery(r *http.Request) url.Values {
	if q, ok := Get(r, "query"); ok {
		return q.(url.Values)
	}

	return nil
}

func QueryParser() Middleware {
	return func(w http.ResponseWriter, r *http.Request) {
		Set(r, "query", r.URL.Query())
		Continue(r)
	}
}
