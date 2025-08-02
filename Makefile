IMAGE_NAME=cloud-server
TAG=latest
PORT=8080

.PHONY: build run build-run

build:
	docker build -t $(IMAGE_NAME):$(TAG) .

run:
	docker run --env-file .env -v /Users/rockon/Storage:/Users/rockon/Storage -p 8080:8080 cloud-server:latest

build-run: build run