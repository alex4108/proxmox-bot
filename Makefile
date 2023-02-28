PHONY: .build .docker .docker-release .test .clean

ifndef PROXMOX_BOT_TAG
override PROXMOX_BOT_TAG = proxmox-bot
endif

.DEFAULT_GOAL := build

clean:
	rm -rf ./bin/

build:
	go mod tidy
	go build -o ./bin/proxmox-bot

docker: build
	docker build -t $(PROXMOX_BOT_TAG) .

docker-release: build
	docker buildx build --platform linux/amd64 -t alex4108/proxmox-bot:$(PROXMOX_BOT_TAG) --push .

test: docker
	docker run --rm -e CI=true $(PROXMOX_BOT_TAG)
