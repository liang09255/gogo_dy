VERSION = v1.0.1

docker_build:
	docker build -t liang09255/gogo_dy:$(VERSION) .

docker_push:
	docker push liang09255/gogo_dy:$(VERSION)