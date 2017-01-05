all: 
	go build -o acro

install:
	go install

test:
	go test -v

run:
	go build -o acro && ./acro
