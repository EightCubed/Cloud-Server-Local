IMAGE_NAME=cloud-server
TAG=latest
PORT=8080

include .env

USER_ID=$(shell id -u)
GROUP_ID=1003

.PHONY: build run build-run

build:
	docker build -t $(IMAGE_NAME):$(TAG) .

run:
	docker run --env-file .env \
		-v $(FILE_STORAGE_PATH)$(FILE_STORAGE_DIRECTORY):$(FILE_STORAGE_PATH)$(FILE_STORAGE_DIRECTORY) \
		--user $(USER_ID):$(GROUP_ID) \
		-p $(PORT):$(PORT) \
		$(IMAGE_NAME):$(TAG)

build-run: build run