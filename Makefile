.PHONY: build sync

DESTINATION=media:media/showsync
BINARY=/tmp/showsync
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS=-ldflags "-X github.com/seamusv/show-sync/cmd.Build=$(BUILD) -s -w"

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(LDFLAGS) main.go
	ls -lh $(BINARY)

sync:
	scp $(BINARY) $(DESTINATION)