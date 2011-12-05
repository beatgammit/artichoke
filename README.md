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

> func RunLocal(port)

* port- int

Convenience for running a server on localhost.

Middleware
----------

This is what middleware looks like:

    func(http.ResponseWriter, *http.Request, map[string]interface{}) bool

The last parameter is for arbitrary data passing. Middleware must define what keys they use and should provide a way to change the default for maximum flexibility. Examine the packaged middleware for examples on style.
