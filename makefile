
build: generate install

generate:
	go generate ./internal/static/...

install:
	go install ./cmd/...
