GIT_TAG := $(shell git describe --abbrev=0 --tags)
PATH_VERSION := $(shell git rev-list $(GIT_TAG).. --count)
GIT_VERSION := "$(GIT_TAG).$(PATH_VERSION)"

.PHONY: build, version, push

version:
	@echo $(GIT_VERSION)

build:
	go build -o ./dist/conveycode.exe \
		-ldflags "-X main.VERSION=$(GIT_VERSION)" \
		./cmd/

push:
	git push
	git push --tags