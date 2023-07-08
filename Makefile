.PHONY: build clean

build: gomodgen
	export GO111MODULE=on
	go build -ldflags="-s -w" -o bin/imaiot main.go

clean:
	rm -rf ./bin