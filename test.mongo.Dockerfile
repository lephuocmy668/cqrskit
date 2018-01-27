FROM mongrel:0.0.1

RUN apk add --no-cache go && rm -rf /var/cache/apk/*

COPY . ./cqrskit
WORKDIR /cqrskit

ENV API_MONGO_TEST_HOST 0.0.0.0:27017
ENV API_MONGO_TEST_DB test_db
ENV API_MONGO_TEST_AUTHDB test_db

CMD ["/bin/bootmgo --fork && go test -v ./repositories/mgorp/..."]