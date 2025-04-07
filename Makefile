.PHONY: build
build:
		go build -o crypto_api/main.go

.DEFAULT_GOAL := build