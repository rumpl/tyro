all: cli plugin-mkdir plugin-mod

cli:
	@go build .

plugin-mkdir:
	@go build -o pls/mkdir ./pkg/plugins/mkdir/mkdir.go

plugin-mod:
	@go build -o pls/mod ./pkg/plugins/mod/mod.go

.PHONY: cli
