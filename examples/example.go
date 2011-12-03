package main

import (
	"artichoke"
	"http"
)

func helloWare(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	w.Write([]byte("Hello world"))
	w.Write([]byte(""))
	return true
}

func router(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	if r.URL.Path == "/" {
		return helloWare(w, r, m)
	}

	return false
}

func main() {
	server := artichoke.New(nil, router)
	server.RunLocal(3345)
}
