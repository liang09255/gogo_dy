VERSION = v2.0.0

build_all:
	make build_base
	make build_api
	make build_user
	make build_chat
	make build_video

build_base:
	docker build -t gogo_base:$(VERSION) --build-arg GOGO_VERSION=$(VERSION) .

build_api:
	docker build -t gogo_api:$(VERSION) --build-arg GOGO_VERSION=$(VERSION) ./api

build_user:
	docker build -t gogo_user:$(VERSION) --build-arg GOGO_VERSION=$(VERSION) ./user

build_chat:
	docker build -t gogo_chat:$(VERSION) --build-arg GOGO_VERSION=$(VERSION) ./chat

build_video:
	docker build -t gogo_video:$(VERSION) --build-arg GOGO_VERSION=$(VERSION) ./video

tidy_all:
	cd ./api && go mod tidy
	cd ./user && go mod tidy
	cd ./video && go mod tidy
	cd ./chat && go mod tidy
	cd ./common && go mod tidy
