VERSION=0.1.0
NOVENDOR=$(shell glide novendor)

bundle:
	glide install

all-build: memds-build memds-cli-build

memds-build: cmd/memds/main.go memds/*.go
	go build -ldflags "-X main.version=${VERSION}" -o bin/memds cmd/memds/main.go

memds-cli-build:
	go build -ldflags "-X main.version=${VERSION}" -o bin/memds-cli cmd/memds-cli/main.go

fmt:
	@echo $(NOVENDOR) | xargs go fmt

test:
	go test -v -cover $(NOVENDOR)

bench:
	go test -bench . $(NOVENDOR)

vet:
	@echo $(NOVENDOR) | xargs go vet

lint:
	golint memds

clean:
	rm -rf bin/*
