clean:
	go clean -r -cache -testcache -modcache
.PHONY: clean

tidy:
	go mod tidy -v
.PHONY: tidy

build:
	go build -buildvcs=false -ldflags '-extldflags "-static"' -o ./bin/ ./ingest/cmd
.PHONY: build

build-clean: clean build
.PHONY: build-clean

test:
	go test -race -failfast -v ./...
.PHONY: test
