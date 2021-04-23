build:
	go build -o main -ldflags="-X 'main.version=dev' -X 'main.commit=$(shell git rev-parse --short HEAD)' -X 'main.builtDate=$(shell date)' -X 'main.builtBy=$(shell hostname)'"

text-proxyjump: build
	./main --verbose --columns host,port,username,jump,nojumps --file test/proxy

test-update: build
	./main --verbose --update

