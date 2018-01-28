FROM influx6/gomongrel:latest

RUN mkdir -p /go/src/github.com/gokit
COPY . /go/src/github.com/gokit/cqrskit

ENV GOPATH /go

WORKDIR /go/src/github.com/gokit/cqrskit
RUN go get ./...
RUN chmod +x -R ./scripts

ENV MONGO_HOST 0.0.0.0:27017
ENV MONGO_USER "test_user"
ENV MONGO_PASSWORD "123456"
ENV MONGO_DB test_db
ENV MONGO_AUTHDB test_db

CMD ["./scripts/mgo.sh"]