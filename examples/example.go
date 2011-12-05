package main

import (
	"artichoke"
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
	params := m["params"].(map[string]string)
	w.Write([]byte("Hello " + params["first"] + " " + params["last"]));
	w.Write([]byte(""))
	return true;
}

func genRoutes() []artichoke.Route {
	ret := []artichoke.Route{
		artichoke.Route{
			Method: "GET",
			Pattern: "/greet/:first/?:last?",
			Handler: handler,
		},
	}

	return ret
}

func logger(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
	fmt.Println("Method: " + r.Method)
	fmt.Println("URL: " + r.URL.Raw)
	fmt.Println("")
	return false
}

func main() {
	server := artichoke.New(nil, logger, artichoke.Router(genRoutes()), artichoke.Static("./public"))
	server.RunLocal(3345)
}
