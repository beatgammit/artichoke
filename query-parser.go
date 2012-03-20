package artichoke

import (
	"net/http"
	"net/url"
)

func (d *Data) GetQuery() (url.Values) {
	if q, ok := d.raw["query"]; ok {
		return q.(url.Values)
	}

	return nil
}

func QueryParser() Middleware {
	return func(w http.ResponseWriter, r *http.Request, m *Data) bool {
		m.raw["query"] = r.URL.Query()
		return false
	}
}
