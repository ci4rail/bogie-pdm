all: container
VERSION=a2

container:
	docker buildx build --build-arg VERSION=${VERSION} -f cmd/bogie-edge/Dockerfile \
		--platform linux/arm64 --push -t ci4rail/bogie-edge:${VERSION} .

.PHONY: all container
