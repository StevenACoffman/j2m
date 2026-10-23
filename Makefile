.DEFAULT_GOAL := easy
.PHONY: install clean all easy test fmt test bin/j2m bin/m2j

bin/j2m:
	go build -o bin/j2m ./cmd/j2m

bin/m2j:
	go build -o bin/m2j ./cmd/m2j

all: fmt bin/j2m bin/m2j test

test:
	go test -cover ./...


install: bin/j2m bin/m2j
	cp bin/* ~/bin

clean:
	rm -f bin/*

easy: all

fmt:
	go fmt ./...
