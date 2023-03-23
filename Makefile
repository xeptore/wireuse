clean:
	go clean -r -cache -testcache -modcache
.PHONY: clean

tidy:
	go mod tidy -v -x
.PHONY: tidy

build:
	rm -rfv ./bin
	mkdir -vp ./bin
	go build -trimpath -buildvcs=false -ldflags '-extldflags "-static" -s -w -buildid=' -o ./bin/ingest ./ingest/cmd
.PHONY: build

build-clean: clean build
.PHONY: build-clean

test:
	go test -trimpath -buildvcs=false -ldflags '-extldflags "-static" -s -w -buildid=' -race -failfast -vet=all -covermode=atomic -coverprofile=coverage.out -v ./...
.PHONY: test

gen:
	$(MAKE) -C ./ingest
.PHONY: gen
