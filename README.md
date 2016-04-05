Deprecation Notice
==================

I am no longer maintaining this and there are known bugs that won't be fixed. I recommend switching to [gin](https://github.com/gin-gonic/gin), as I've seen better stability and lower resource usage along with better structure..

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

> func Run(host, port)

* host- string
* port- int

Runs the server on the specified port and host.

> func RunTLS(host, port, certFile, keyFile)

* host- string
* port- int
* certFile- string; path to TLS certificate
* keyFile- string; path to private key

Runs the server on the specified port and host using the certificate and private key for TLS sessions.

Middleware
----------

This is what middleware looks like:

    func(http.ResponseWriter, *http.Request, Data) bool

The last parameter is for arbitrary data passing. Middleware attach functions to this struct. Internally, there is a hidden `raw` property (`map[string]interface{}`) that middleware should attach their data to. Examine the packaged middleware for examples on style.

**BasicAuth**

Performs basic HTTP authentication.

    func BasicAuth(auth, required) Middleware

* auth- map[string]string; maps a username to a password
* required- bool; if true, sends a 401 and ends the response

*BasicAuth* will make a new key and attach a `GetAuth()` function to the Data parameter, which returns a struct with the following properties:

* `User`- string; user name provided
* `Pass`- string; password provided
* `Authenticated`- bool; whether the user was successfully authenticated

If no credentials were passed in, *auth* will be `nil`.

**BodyParser**

Attempts to parse the body as defined by the content-type.

    func BodyParser(maxMemory) Middleware

* maxMemory- int64; maximum memory, in bytes, the request can use for caching (for multipart)
  * Once the max memory has been used, files will be written to the filesystem

*BodyParser* will make an attempt to decode the request body depending on the content type.
Each request will populate a Body struct, which has the following properties:

* `Data`- representation of the body that was parsed; `interface{}`
* `Raw`- raw body, if available; []byte
* `Err`- error during parsing, if any

For example:

* `application/json`- parses body as JSON and stores result of json.Unmarshal in Body.Data
* `application/x-www-form-encoded`- uses the built-in r.ParseForm(); result is stored in r.Form
* `multipart/form-data`- uses the built-in r.ParseMultipartForm; result is stored in r.Form
* default- full-text of the body is stored in Data.Data as a string

**Router**

Basic request router. There are two ways to instantiate:

*Static Routes*

This is the simplest way to create a router, especially if the routes will never change:

    func StaticRouter(...routes) Middleware

*Dynamic Routes*

To have a dynamic router, create a new `Router`:

    func NewRouter(...Routes) Router

These public functions will be useful:

    func (r Router) Add(...Routes)
    func (r Router) Middleware() Middleware

*Router*

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

Router attaches a `GetParams()` function to the Data parameter, which returns a `Params` struct, which has a `Get(key string)` function that gets the value for a specified key, which is also a string.

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

Simple static file handler that ends the response only if the file exists, or there is an error reading the file. This means that it is possible for subsequent Middleware to get called after Static.

    func Static(root) Middleware

* root- string; directory root from current directory

**QueryParser**

Uses the built-in url.URL.Query() and attaches a `GetQuery()` function to the Data map that returns a url.Values object (basically a map[string][]string).

    func QueryParser() Middleware

There are no parameters, because it's super simple.
