clean:
	go clean -r -cache -testcache -modcache
.PHONY: clean

tidy:
	go mod tidy -v -x
.PHONY: tidy

build:
	rm -rfv ./bin
	mkdir -vp ./bin
	go build -trimpath -buildvcs=false -ldflags '-extldflags "-static" -s -w -buildid=' -o ./bin/ingest_$$(go env GOOS)_$$(go env GOARCH) ./ingest/cmd
.PHONY: build

build-clean: clean build
.PHONY: build-clean

test:
	go test -trimpath -buildvcs=false -ldflags '-extldflags "-static" -s -w -buildid=' -race -failfast -vet=all -v ./...
.PHONY: test
