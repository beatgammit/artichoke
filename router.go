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

func (d *Data) GetParams() *Params {
	if p, ok := d.raw["params"]; ok {
		return p.(*Params)
	}
	return nil
}

func Router(routes []Route) Middleware {
	for i, v := range routes {
		if v.Handler == nil {
			panic("Every route must have a func that implements artichoke.Middleware")
		}
		if v.Method == "" {
			panic("Every route must have a method: GET, POST, etc.")
		}

		routes[i].Method = strings.ToUpper(v.Method)

		switch t := v.Pattern.(type) {
		case string:
			pattern := t

			// turn sinatra-style routing into a named submatch
			// for example: '/:root' => '/(?P<root>[^/?#])'
			pattern = varRegex.ReplaceAllString(pattern, "(?P<$1>[^/?#]*)")

			// store the into this Route object
			// go ahead and panic; all panics will occur during debugging anyway
			routes[i].reg = regexp.MustCompile(pattern)

			// grab the variable names from the regex
			// only computed this once
			routes[i].vars = routes[i].reg.SubexpNames()[1:]
		case regexp.Regexp:
			routes[i].reg = &t
		default:
			if _, ok := v.Pattern.(*regexp.Regexp); ok {
				panic("Pattern is not a string or a regexp!")
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request, d *Data) bool {
		for _, v := range routes {
			// use Contains because v.Method could have comma-separated methods
			if !strings.Contains(v.Method, r.Method) && v.Method != "*" {
				continue
			}

			// check if there's a match
			matches := v.reg.FindAllStringSubmatch(r.URL.Path, -1)
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

			d.raw["params"] = NewParams(params)
			if res := v.Handler(w, r, d); res {
				return true
			}
		}
		return false
	}
}
