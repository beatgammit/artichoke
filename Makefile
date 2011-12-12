include $(GOROOT)/src/Make.inc

TARG=artichoke
GOFILES=\
	src/core.go\
	src/middleware/static.go\
	src/middleware/router.go\
	src/middleware/basic-auth.go\
	src/middleware/query-parser.go\
	src/middleware/body-parser.go\

# makes a package
include $(GOROOT)/src/Make.pkg
