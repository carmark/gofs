export GOPATH:=$(abs_top_srcdir)/Godeps/_workspace:$(GOPATH)

all-local: build
clean-local:
	-rm -f gofs
install-exec-local: 
	$(INSTALL_PROGRAM) gofs $(bindir)

build:
	go build gofs.go
