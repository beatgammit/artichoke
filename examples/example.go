package main

import (
	"../"
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
	params := artichoke.GetParams(m)
	w.Write([]byte("Hello " + params.Get("first") + " " + params.Get("last")));
	w.Write([]byte(""))
	return true;
}

func genRoutes() []*artichoke.Route {
	ret := []*artichoke.Route{
		&artichoke.Route{
			Method: "GET",
			Pattern: "/greet/:first/?:last?",
			Handler: handler,
		},
	}

	return ret
}

func logger(w http.ResponseWriter, r *http.Request, m artichoke.Data) bool {
	fmt.Println("Method:", r.Method)
	fmt.Println("URL:", r.URL.Path)

	// auth can  be nil if no authentication data was passed in
	if auth := artichoke.GetAuth(m); auth != nil {
		fmt.Println("User:", auth.User)
		fmt.Println("Password:", auth.Pass)
		fmt.Println("Authenticated:", auth.Authenticated)
	} else {
		fmt.Println("No authentication data provided")
	}

	fmt.Println("Query:")
	for k, vals := range artichoke.GetQuery(m) {
		for _, v := range vals {
			fmt.Println("  " + k + " : " + v)
		}
	}

	if body := artichoke.GetBody(m); body != nil {
		fmt.Println("Body:")
		fmt.Println("  " + string(body.Raw))
	}

	fmt.Println()
	return false
}

func main() {
	server := artichoke.New(nil,
			artichoke.BasicAuth(map[string]string{"jack": "johnson"}, false),
			artichoke.QueryParser(),
			artichoke.BodyParser(1024 * 10),
			logger,
			artichoke.StaticRouter(genRoutes()...),
			artichoke.Static("./public"),
		)
	server.Run("localhost", 3345)
}
