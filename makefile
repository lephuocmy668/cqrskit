
build: generate install

generate:
	go generate ./internal/static/...

install:
	go install ./cmd/...

buildMGO:
	docker build -t cqrskit-mgo-test -f test.mongo.Dockerfile ./

testMGO: buildMGO
	docker run -it --rm cqrskit-mgo-test
