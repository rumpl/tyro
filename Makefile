PROTOS=$(shell find protos -name \*.proto)

all: cli plugin-mkdir plugin-mod protos

cli:
	go build .

plugin-mkdir:
	go build -o pls/mkdir ./pkg/plugins/mkdir/mkdir.go

plugin-mod:
	go build -o pls/mod ./pkg/plugins/mod/mod.go

protos:
	protoc -I. --go_out=plugins=grpc:. ${PROTOS}

.PHONY: cli plugin-mkdir plugin-mod protos
