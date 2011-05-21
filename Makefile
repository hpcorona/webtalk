
include $(GOROOT)/src/Make.inc

TARG = webtalk
GOFILES = \
	interop.go \
	frame.go \
	talkers.go \
	management.go \
	main.go

include $(GOROOT)/src/Make.cmd

run: all
	./webtalk
