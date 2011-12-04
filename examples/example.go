package main

import (
	"artichoke"
	"http"
	"regexp"
)

func router(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
	reg, _ := regexp.Compile("/greet/([A-Za-z0-9]*)/?([A-Za-z0-9]*)?")
	if reg.MatchString(r.URL.Path) {
		var s string
		matches := reg.FindStringSubmatch(r.URL.Path);
		for _, match := range(matches[1:]) {
			s += " " + match
		}
		w.Write([]byte("Hello" + s))
		w.Write([]byte(""))
		return true
	}

	return false
}

func main() {
	server := artichoke.New(nil, router, artichoke.Static("./public"))
	server.RunLocal(3345)
}
