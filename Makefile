VERSION = v2.0.0

build_all:
	make build_base
	make build_api

build_base:
	docker build -t base:$(VERISON) .

build_api:
	docker build -t api:$(VERSION) ./api

tidy_all:
	cd ./api && go mod tidy
	cd ./user && go mod tidy
	cd ./video && go mod tidy
	cd ./chat && go mod tidy
	cd ./common && go mod tidy
