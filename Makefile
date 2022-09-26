.DEFAULT_GOAL := build-n-install

build-n-install:
	make build
	make install

build:
	go build -o bin/ipstat -ldflags "-s -w"

docker-build:
	docker build -t ipstat .

docker-run:
	docker run --rm -p 3000:3000 ipstat

publish:
	make publish-ghcr

publish-ghcr:
	#make docker-build
	docker tag ipstat:latest docker.pkg.github.com/nathan13888/ipstat/api:latest
	docker push docker.pkg.github.com/nathan13888/ipstat/api:latest

pull-ghcr:
	docker pull docker.pkg.github.com/nathan13888/ipstat/api:latest

test:
	go test -v ./...
