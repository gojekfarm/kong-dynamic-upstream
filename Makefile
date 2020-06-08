.PHONY: all

all: clean test dynamicupstream.so

clean:
	rm *.so

test:
	go test ./...

dynamicupstream.so:
	go build -o dynamicupstream.so -buildmode=plugin handler.go
