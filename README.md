Intro
=====

Similar to Node's connect module, `artichoke` offers a layered approach to web application programming.

Everything is a node in a stack. Starting at the top of the stack (the first function passed in), each function is executed until one ends the response, and the rest of the functions are not executed. If the stack has been traversed and the response has not been ended, a default error handler will end it.

**This project is currently being developed using the weekly builds of Go. This is to be considered highly unstable and the API may change drastically.**

Philosophy
----------

> "Write middleware that do one thing and do it well"

Each middleware should accomplish a single task. This greatly reduces the complexity of applications, which can greatly increase performance.

The provided middleware are meant to be simple. They pass information down the stack, and only end the response when necessary. An example stack might look like this:

1. Authenticate user (end response with error code if not in database)
2. Parse body- best attempt to decode the body (turn JSON string into Go map)
3. Router- handle special API calls (ends the response if route exists)
4. Handle static files (always ends the response with appropriate error codes)

The platform is not meant to be used as a framework by itself, but as a framework to build other frameworks.

API
===

There are only a few important functions that should be called:

> func New(options, middleware)

* options- map[string] interface
* middleware- ...Middleware

This creates a new Server. `options` is not currently being used, but it will eventually have configurable options to set default behavior.

All Middleware will be passed to Use.

**Attached to Server**

> func Use(middleware)

* middleware- ...Middleware

Add middleware. Middleware will be called in the order received. For more information, please see the middleware section below.

> func Run(port, host)

* port- int
* host- string

Runs the server on the specified port and host.

> func RunTLS(port, host, certFile, keyFile)

* port- int
* host- string
* certFile- string; path to TLS certificate
* keyFile- string; path to private key

Runs the server on the specified port and host using the certificate and private key for TLS sessions.

Middleware
----------

This is what middleware looks like:

    func(http.ResponseWriter, *http.Request, map[string]interface{}) bool

The last parameter is for arbitrary data passing. Middleware must define what keys they use and should provide a way to change the default for maximum flexibility. Examine the packaged middleware for examples on style.

**BasicAuth**

Performs basic HTTP authentication.

    func BasicAuth(auth, required) Middleware

* auth- map[string]string; maps a username to a password
* required- bool; if true, sends a 401 and ends the response

*BasicAuth* will make a new key in the Data parameter passed down the stack named *auth* of type map[string]interface{}. The keys of this map will be:

* `user`- string; user name provided
* `pass`- string; password provided
* `authenticated`- bool; whether the user was successfully authenticated

If no credentials were passed in, *auth* will be `nil`.

**BodyParser**

Attempts to parse the body as defined by the content-type.

    func BodyParser(maxMemory) Middleware

* maxMemory- int64; maximum memory, in bytes, the request can use for caching (for multipart)
  * Once the max memory has been used, files will be written to the filesystem

*BodyParser* will make an attempt to decode the request body depending on the content type. For example:

* `application/json`- parses body as JSON and stores result in Data["bodyJson"] and the raw text in Data["body"]
* `application/x-www-form-encoded`- uses the built-in r.ParseForm(); result is stored in r.Form
* `multipart/form-data`- uses the built-in r.ParseMultipartForm; result is stored in r.Form
* default- full-text of the body is stored in Data["body"]

If an error occurs, the error is stored in d["bodyParseError"]. This will be of type os.Error.

**Router**

Basic request router. Give it an array of routes and it will do the routing work for you.

    func Router(routes) Middleware

* routes- []Route; array of Route instances

	type Route struct {
		Method string
		Pattern interface{}
		Handler Middleware
	}

* Method- GET, POST, etc. or * to handle all request types
* Pattern- string or regexp.Regexp; if it's a string, it will be parsed Sinatra style
    * each :identifier maps to a position in the URL
	* ? makes the previous character or group optional
	* this gets parsed into a regexp.Regexp object
* Handler- implementation of Middleware to handle the routed response
    * if more than one pattern matches, they will be treated as individual middleware

*Notes on Pattern:*

Router creates a `params` key in the Data map that is passed down the stack. `params` has the form: `map[string]string`, with keys being identifiers, and values being the matches in the URL

Examples:

* `/:first/:second`
    * '/hello/world'- matched
        * first- hello
        * second- world
    * '/hello'- not matched
    * '/hello/'- not matched
    * '/'- not matched
* `/:first/:second?`
    * '/hello/world'- matched
        * first- hello
        * second- world
    * '/hello'- not matched
    * '/hello/'- matched
        * first- hello
    * '/'- not matched
* `/:first/?:second?`
    * '/hello/world'- matched
        * first- hello
        * second- world
    * '/hello'- matched
        * first- hello
    * '/hello/'- matched
        * first- hello
    * '/'- not matched

It's pretty simple. Be creative! For example, to match file extensions for a file named File (to dynamically respond to requests for different file formats):

    /File.:extension?

**Static**

Simple static file handler that always ends the response.

	func Static(root) Middleware

* root- string; directory root from current directory

**QueryParser**

Uses the built-in url.URL.Query() and populates a `query` key in the Data map that is passed down the stack with a url.Values object as the value (basically a map[string][]string).

	func QueryParser() Middleware

There are no parameters, because it's super simple.
