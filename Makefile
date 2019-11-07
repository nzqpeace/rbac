Binary=rbac

.PHONY: get_deps build all test clean

all:build

get_deps:
	go get ./...

build:get_deps
	go build -o ./rbac/${Binary} ./rbac/
	go build -o ./example/example ./example/

test:
	go test -v .
	go test -v ./db
	go test -v ./cache

clean:
	go clean -i -x; rm rbac/${Binary}; rm example/example
