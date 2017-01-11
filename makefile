all: 
	go build

install:
	go install

test:
	go test -v

run:
	go build && ./acromantula
