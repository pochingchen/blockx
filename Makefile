build:
	go build -o ./build/blockx

run: clean build
	./build/blockx

test:
	go test ./...

clean:
	rm -rf ./build