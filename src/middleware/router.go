package artichoke

import (
	"regexp"
	"net/http"
	"strings"
)

// Unicode isn't supported yet in all browsers, so we'll take the least common denominator
const URLChars = "[A-Za-z0-9$_.+!*'\"(),-]"

// :identifier
// this is a full Unicode-compliant regex
var varRegex = regexp.MustCompile("(:[\\p{Lu}\\p{Ll}\\p{Lt}\\p{Lm}\\p{Lo}][\\p{Lu}\\p{Ll}\\p{Lt}\\p{Lm}\\p{Lo}\\p{Nd}]*)")

type Route struct {
	// GET, POST, PUT, etc.
	Method string
	// passed by the client as a string or regexp, but always a regexp when used
	Pattern interface{}
	Handler Middleware
	// variable names by position in regexp match
	// this will be not be nil only if Pattern has at least one variable
	vars []string
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

			// get all of the variable names
			vars := varRegex.FindAllString(pattern, -1)
			for j, val := range vars {
				vars[j] = val[1:]
			}

			routes[i].vars = vars

			// get rid of anything that could mess everything up
			// since urls can have regex characters (like parens, dots, and *),
			// make sure we don't get any weird regexes
			pattern = regexp.QuoteMeta(pattern)

			// replace all variables in the pattern with a regex to match a part of a URL
			pattern = varRegex.ReplaceAllString(pattern, "(" + URLChars + "+)")

			// allow optional parts of the string
			pattern = strings.Replace(pattern, "\\?", "?", -1)

			// store the generated regex into this Route object
			// go ahead and panic; all panics will occur during debugging anyway
			routes[i].Pattern = regexp.MustCompile(pattern)
		case regexp.Regexp:
			routes[i].Pattern = &t
		default:
			if _, ok := v.Pattern.(*regexp.Regexp); ok {
				panic("Pattern is not a string or a regexp!")
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request, d Data) bool {
		for _, v := range routes {
			// use Contains because v.Method could have comma-separated methods
			if !strings.Contains(v.Method, r.Method) && v.Method != "*" {
				continue
			}

			reg := v.Pattern.(*regexp.Regexp)

			// check if there's a match
			matches := reg.FindAllStringSubmatch(r.URL.Raw, -1)
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

			d["params"] = params
			if res := v.Handler(w, r, d); res {
				return true
			}
		}
		return false
	}
}
