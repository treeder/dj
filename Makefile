build: 
	go build -o dj

install:
	go build -o /usr/local/bin/dj

docker: 
	docker build -t treeder/dj:latest .

dep:
	glide install -v

dep-up:
	glide up -v

test:
	./test.sh

release:
	GOOS=linux go build -o dj_linux
	GOOS=darwin go build -o dj_mac
	GOOS=windows go build -o dj.exe
	docker run --rm -v ${PWD}:/go/src/github.com/treeder/dj -w /go/src/github.com/treeder/dj golang:alpine go build -o dj_alpine

.PHONY: install
