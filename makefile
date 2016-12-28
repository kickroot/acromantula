all: 
	go build -o acro

install:
	go install

test:
	go test -v
