GIT_TAG := $(shell git describe --abbrev=0 --tags)

.PHONY: build, version

version:
	@echo $(GIT_TAG)

build:
	go build -o ./dist/conveycode.exe \
		-ldflags "-X main.VERSION=$(GIT_TAG)" \
		./cmd/

