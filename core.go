package artichoke

import (
	"fmt"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
)

var errors = map[int]string{
	http.StatusNotFound:            fmt.Sprintf("<h1>Error %d: Not Found</h1><br /><br />The page or resource requested could not be found. If this was a link or worked previously, please notify your webmaster.", http.StatusNotFound),
	http.StatusInternalServerError: fmt.Sprintf("<h1>Error %d: Internal Server Error</h1><br /><br />An internal server error prevented execution of this request. Please notify the webmaster.", http.StatusInternalServerError),
}

type Data interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}

type data struct {
	raw map[string]interface{}
}

func (d *data) Get(key string) (interface{}, bool) {
  i, ok := d.raw[key]
  return i, ok
}

func (d *data) Set(key string, val interface{}) {
	d.raw[key] = val
}

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
	data := new(data)
	data.raw = make(map[string]interface{})

	for _, fn := range(s.middleware) {
		if fn(w, r, data) == true {
			return
		}
	}

	status := http.StatusNotFound
	fmt.Println("No handler for this request:")
	fmt.Printf("  Method: %s\n", r.Method)
	fmt.Printf("  URL: %s\n", r.URL.Path)
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

func (s *Server) Run(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		panic(e)
	}

	s.l = l
	srv := &http.Server{Addr: addr, Handler: s}

	fmt.Println("Server starting on port:", port)
	srv.Serve(s.l)
}

func (s *Server) RunTLS(host string, port int, certFile string, keyFile string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{NextProtos: []string{"http/1.1"}}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		panic(err)
	}

	s.l = l
	srv := &http.Server{Addr: addr, Handler: s, TLSConfig: config}

	fmt.Println("Secure server starting on port:", port)
	srv.Serve(s.l)
}

func (s *Server) Stop() {
	s.l.Close()
}
