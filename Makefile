all: docker-arm64
VERSION=a1

docker-arm64:
	docker buildx build --build-arg VER=${VERSION} -f cmd/bogie-edge/Dockerfile \
		--platform linux/arm64 --push -t ci4rail/bogie-edge:${VERSION} .

.PHONY: all docker-arm64
