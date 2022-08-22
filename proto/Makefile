all: go/metrics/v1/metrics.go

prepare:
	mkdir -p go/
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go/metrics/v1/metrics.go: metrics/v1/metrics.proto prepare
	protoc -I=./metrics ./metrics/v1/metrics.proto --go_out=go/

.PHONY: all prepare
