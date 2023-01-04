SRC_DIR=.
all: ${SRC_DIR}/go/metrics/v1/metrics.go \
 ${SRC_DIR}/go/bogie/v1/bogie.go \
 ${SRC_DIR}/python/bogie/v1/bogie_pb2.py \
 ${SRC_DIR}/python/metrics/v1/metrics_pb2.py

prepare:
	mkdir -p go/ python/
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

${SRC_DIR}/go/metrics/v1/metrics.go: ${SRC_DIR}/metrics/v1/metrics.proto prepare
	protoc -I=${SRC_DIR} ${SRC_DIR}/metrics/v1/metrics.proto --go_out=go/

${SRC_DIR}/go/bogie/v1/bogie.go: ${SRC_DIR}/bogie/v1/bogie.proto prepare
	protoc -I=${SRC_DIR} ${SRC_DIR}/bogie/v1/bogie.proto --go_out=go/

${SRC_DIR}/python/bogie/v1/bogie_pb2.py: ${SRC_DIR}/bogie/v1/bogie.proto prepare
	protoc -I=${SRC_DIR} ${SRC_DIR}/bogie/v1/bogie.proto --python_out=python/
	cp ${SRC_DIR}/python/bogie/v1/bogie_pb2.py  ../analytics/proto

${SRC_DIR}/python/metrics/v1/metrics_pb2.py: ${SRC_DIR}/metrics/v1/metrics.proto prepare
	protoc -I=${SRC_DIR} ${SRC_DIR}/metrics/v1/metrics.proto --python_out=python/
	cp ${SRC_DIR}/python/metrics/v1/metrics_pb2.py  ../analytics/proto


.PHONY: all prepare
