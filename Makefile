include $(GOROOT)/src/Make.inc

TARG=artichoke
GOFILES=\
	src/core.go\

# makes a package
include $(GOROOT)/src/Make.pkg
