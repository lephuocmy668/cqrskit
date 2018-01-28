#!/usr/bin/env sh

env AUTH=true /bin/bootmgo --fork

go test -v ./repositories/mgorp/...