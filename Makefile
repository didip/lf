# Makefile for

all: lf

lf: FORCE
	go build -o lf -a cmd/lf/*.go

clean:	FORCE
	rm -rf lf lf-db-test

godeps:	FORCE
	go get -u github.com/NYTimes/gziphandler
	go get -u github.com/codahale/rfc6979
	go get -u github.com/tidwall/pretty
	go get -u golang.org/x/crypto/ed25519
	go get -u golang.org/x/crypto/sha3
	go get -u gopkg.in/kothar/brotli-go.v0/enc
	go get -u gopkg.in/kothar/brotli-go.v0/dec

FORCE:
