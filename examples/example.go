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
	fmt.Println("Method:", r.Method)
	fmt.Println("URL:", r.URL.Raw)

	// auth can  be nil if no authentication data was passed in
	if m["auth"] != nil {
		auth := m["auth"].(map[string]interface{})

		fmt.Println("User:", auth["user"].(string))
		fmt.Println("Password:", auth["pass"].(string))
		fmt.Println("Authenticated:", auth["authenticated"].(bool))
	} else {
		fmt.Println("No authentication data provided")
	}

	fmt.Println()
	return false
}

func main() {
	server := artichoke.New(nil,
			artichoke.BasicAuth(map[string]string{"jack": "johnson"}, false),
			logger,
			artichoke.Router(genRoutes()),
			artichoke.Static("./public"),
		)
	server.RunLocal(3345)
}
