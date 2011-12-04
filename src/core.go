package artichoke

import (
	"http"
	"net"
	"fmt"
	"strconv"
	"container/vector"
)

var errors = map[int] string {
	http.StatusNotFound: fmt.Sprintf("<h1>Error %d: Not Found</h1><br /><br />The page or resource requested could not be found. If this was a link or worked previously, please notify your webmaster.", http.StatusNotFound),
	http.StatusInternalServerError: fmt.Sprintf("<h1>Error %d: Internal Server Error</h1><br /><br />An internal server error prevented execution of this request. Please notify the webmaster.", http.StatusInternalServerError),
}

type Data map[string]interface{}

// once a middleware returns true, no more middleware will be executed
//
// the last parameter is a general-purpose map passed to each middleware
// middleware can use this to pass arbitrary data down the stack
type Middleware func(http.ResponseWriter, *http.Request, Data) bool

type Server struct {
	handler func (http.ResponseWriter, *http.Request)
	middleware vector.Vector // allows middleware to be added after the New
	l net.Listener
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
	for _, fn := range fns {
		s.middleware.Push(fn)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data Data
	for i := 0; i < s.middleware.Len(); i++ {
		fn := s.middleware.At(i).(Middleware)
		if fn(w, r, data) == true {
			return
		}
	}

	status := http.StatusNotFound
	fmt.Println("No handler for this request:")
	fmt.Printf("  Method: %s\n", r.Method)
	fmt.Printf("  URL: %s\n", r.RawURL)
	fmt.Println("  Headers:")
	for k, v := range(r.Header) {
		fmt.Printf("    %s: %s\n", k, v)
	}
	fmt.Println("")

	resp := errors[status];
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
	mux := http.NewServeMux()
	mux.Handle("/", s)

	l, err := net.Listen("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println(err.String())
	}
	s.l = l

	fmt.Println("Starting server on port: " + strconv.Itoa(port))
	http.Serve(s.l, mux)
}

func (s *Server) RunLocal(port int) {
	s.Run(port, "localhost")
}
