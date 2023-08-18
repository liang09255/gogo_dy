VERSION = v1.0.1

docker_build:
	docker build -t liang09255/gogo_dy:$(VERSION) .

docker_push:
	docker push liang09255/gogo_dy:$(VERSION)

tidy_all:
	cd ./api && go mod tidy
	cd ./user && go mod tidy
	cd ./video && go mod tidy
	cd ./chat && go mod tidy
	cd ./common && go mod tidy
