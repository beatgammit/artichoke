package middleware

import (
	"github.com/beatgammit/artichoke"
	"net/http"
	"net/url"
)

func GetQuery(d artichoke.Data) url.Values {
	if q, ok := d.Get("query"); ok {
		return q.(url.Values)
	}

	return nil
}

func QueryParser() artichoke.Middleware {
	return func(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
		m.Set("query", r.URL.Query())
		return false
	}
}
