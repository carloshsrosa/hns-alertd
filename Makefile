.PHONY: build test lint fmt bench pre-commit setup clean

build:
	go build -o bin/alertd ./cmd/alertd

# Race detector is non-negotiable. We had a real race in Subsystem.Stats() in
# week 1 that only -race caught; CI runs the same command for that reason.
test:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

bench:
	go test -bench=. -benchmem ./...

pre-commit: fmt lint test

setup:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit

clean:
	rm -rf bin/
