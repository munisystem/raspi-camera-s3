PKG_NAME=$(shell basename `pwd`)

default: build

install:
	go get -t -v ./...

.PHONY: install
build:
	GOOS=linux GOARCH=arm go build -o ./bin/$(PKG_NAME)
