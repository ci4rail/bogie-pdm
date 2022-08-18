all: docker-arm64

docker-arm64:
	docker buildx build -f cmd/gps-bms/Dockerfile --platform linux/arm64 --push -t ci4rail/gps-bms .

.PHONY: all docker-arm64
