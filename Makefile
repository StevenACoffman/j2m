.DEFAULT_GOAL := easy
.PHONY: install clean all easy

bin/j2m:
	go build -o bin/j2m cmd/j2m.go

all: bin/j2m

install: bin/j2m
	cp bin/* ~/bin

clean:
	rm -f bin/*

easy: all