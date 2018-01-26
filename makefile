
build: generate install

generate:
	go generate ./static/...

install:
	go install
