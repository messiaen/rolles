VERSION = 0.0.1

default: clean build

clean:
	@rm -f rolles

rolles:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o rolles github.com/messiaen/rolles/cmd/rolles

build: rolles

# Docker
IMAGE_NAME = "messiaen/rolles"

.PHONY: build_image
build_image: clean build
	@docker build --no-cache -t $(IMAGE_NAME):$(VERSION) .
