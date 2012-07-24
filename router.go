package artichoke

import (
	"regexp"
	"net/http"
	"strings"
)

// :identifier
// this is a full Unicode-compliant regex
var varRegex = regexp.MustCompile(":([\\p{Lu}\\p{Ll}\\p{Lt}\\p{Lm}\\p{Lo}][\\p{Lu}\\p{Ll}\\p{Lt}\\p{Lm}\\p{Lo}\\p{Nd}]*)")

type Route struct {
	// GET, POST, PUT, etc.
	Method string
	// passed by the client as a string or regexp, but always a regexp when used
	Pattern interface{}
	Handler Middleware
	// variable names by position in regexp match
	// this will be not be nil only if Pattern has at least one variable
	vars []string

	// the regex from the user or the Regexp generated
	reg *regexp.Regexp
}

// helper that takes user input and gives it meaning
func prepRoute(r *Route) {
	if r.Handler == nil {
		panic("Every route must have a func that implements artichoke.Middleware")
	}
	if r.Method == "" {
		panic("Every route must have a method: GET, POST, etc.")
	}

	r.Method = strings.ToUpper(r.Method)

	switch t := r.Pattern.(type) {
	case string:
		pattern := t

		// turn sinatra-style routing into a named submatch
		// for example: '/:root' => '/(?P<root>[^/?#])'
		pattern = varRegex.ReplaceAllString(pattern, "(?P<$1>[^/?#]*)")

		// tack on an ending anchor; user must account for it
		if pattern[len(pattern)-1] != '$' {
			pattern += "$"
		}

		// store the into this Route object
		// go ahead and panic; all panics will occur during debugging anyway
		r.reg = regexp.MustCompile(pattern)

	case regexp.Regexp:
		r.reg = &t

	case *regexp.Regexp:
		r.reg = t

	default:
		panic("Pattern is not a string or a regexp!")
	}

	// grab the variable names from the regex
	// only compute this once
	r.vars = r.reg.SubexpNames()[1:]
}

type Router interface {
	Add(...*Route)
	Remove(...*Route)
	Middleware() Middleware
}

type router struct {
	routes []*Route
	sem chan bool
}

func NewRouter(routes ...*Route) Router {
	r := new(router)
	r.sem = make(chan bool, 1)
	r.Add(routes...)

	return r
}

func (r *router) Add(routes ...*Route) {
	for _, r := range routes {
		prepRoute(r)
	}

	// make sure nobody can mess with routes while we're modifying it
	defer func() {
		<-r.sem
	}()
	r.sem <- true

	r.routes = append(r.routes, routes...)
}

func (r *router) Remove(routes ...*Route) {
	var keep []*Route

	// make sure nobody can mess with routes while we're modifying it
	defer func() {
		<-r.sem
	}()
	r.sem <- true

	for _, route := range routes {
		keep = append(keep, route)
	}

	r.routes = keep
}

// returns a closure with access to the router
func (r *router) Middleware() Middleware {
	return func(w http.ResponseWriter, req *http.Request, d Data) bool {
		for _, v := range r.routes {
			// use Contains because v.Method could have comma-separated methods
			if !strings.Contains(v.Method, req.Method) && v.Method != "*" {
				continue
			}

			// check if there's a match
			matches := v.reg.FindAllStringSubmatch(req.URL.Path, -1)
			if matches == nil {
				continue
			}

			// params is in this scope so there is no cross-over of keys between routes
			params := make(map[string]string)

			for _, m := range matches {
				// skip the full string match
				for i, val := range m[1:] {
					params[v.vars[i]] = val
				}
			}

			d.Set("params", NewParams(params))
			if res := v.Handler(w, req, d); res {
				return true
			}
		}
		return false
	}
}

type Params struct {
	raw map[string]string
}

func NewParams(raw map[string]string) *Params {
	p := new(Params)
	p.raw = raw

	return p
}

func (p *Params) Get(key string) string {
	return p.raw[key]
}

func GetParams(d Data) *Params {
	if p, ok := d.Get("params"); ok {
		return p.(*Params)
	}
	return nil
}

func StaticRouter(routes ...*Route) Middleware {
	router := NewRouter(routes...)
	return router.Middleware()
}
