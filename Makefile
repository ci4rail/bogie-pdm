all: container
ARCH?=arm64
CONTAINER_ARCHS?=linux/arm64,linux/amd64
VERSION?=$(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/ci4rail/bogie-pdm/internal/version.Version=${VERSION}"

container:
	docker buildx build --build-arg VERSION=${VERSION} -f cmd/bogie-edge/Dockerfile \
		--platform {CONTAINER_ARCHS} --push -t ci4rail/bogie-edge:${VERSION} .

bogie-edge-static:
	GOOS=linux GOARCH=${ARCH} go build $(GO_LDFLAGS) -o ./bin/bogie-edge-static ./cmd/bogie-edge/main.go


.PHONY: all container bogie-edge-static
