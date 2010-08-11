include $(GOROOT)/src/Make.$(GOARCH)

TARG=ftp
GOFILES=\
	client.go\

include $(GOROOT)/src/Make.pkg
