package artichoke

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

var errors = map[int]string{
	http.StatusNotFound:            fmt.Sprintf("<h1>Error %d: Not Found</h1><br /><br />The page or resource requested could not be found. If this was a link or worked previously, please notify your webmaster.", http.StatusNotFound),
	http.StatusInternalServerError: fmt.Sprintf("<h1>Error %d: Internal Server Error</h1><br /><br />An internal server error prevented execution of this request. Please notify the webmaster.", http.StatusInternalServerError),
}

type Data map[string]interface{}

// once a middleware returns true, no more middleware will be executed
//
// the last parameter is a general-purpose map passed to each middleware
// middleware can use this to pass arbitrary data down the stack
type Middleware func(http.ResponseWriter, *http.Request, Data) bool

type Server struct {
	handler    func(http.ResponseWriter, *http.Request)
	middleware []Middleware
	l          net.Listener
	// for TLS connections
	certFile   string
	keyFile   string
}

var server Server

// create a new server with the options provided
// the first parameter specifies options to control behavior of the server
// any other parameters are just passed to Use for convenience
func New(options map[string]interface{}, fns ...Middleware) *Server {
	s := Server{}
	s.Use(fns...)

	return &s
}

// Adds any number of middleware
// fns is any number of functions that act as middleware
// they will be called order on every request
func (s *Server) Use(fns ...Middleware) {
	s.middleware = append(s.middleware, fns...)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := make(Data)
	for _, fn := range(s.middleware) {
		if fn(w, r, data) == true {
			return
		}
	}

	status := http.StatusNotFound
	fmt.Println("No handler for this request:")
	fmt.Printf("  Method: %s\n", r.Method)
	fmt.Printf("  URL: %s\n", r.URL.Raw)
	fmt.Println("  Headers:")
	for k, v := range r.Header {
		fmt.Printf("    %s: %s\n", k, v)
	}
	fmt.Println("")

	resp := errors[status]
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(status)

	// for HEAD requests, do everything except the body
	if r.Method == "HEAD" {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(resp))
}

func (s *Server) Run(port int, host string) {
	fmt.Println("Server starting on port:", port)

	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), s)
}

func (s *Server) RunTLS(port int, host string, certFile string, keyFile string) {
	fmt.Println("Secure server starting on port:", port)

	http.ListenAndServeTLS(fmt.Sprintf("%s:%d", host, port), certFile, keyFile, s)
}
