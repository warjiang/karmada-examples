.PHONY: run apply check build

run:
	go run .

apply:
	go run . --apply

check:
	gofmt -w .
	go test ./...
	go vet ./...

build:
	go build -o ark-litellm-sync .
