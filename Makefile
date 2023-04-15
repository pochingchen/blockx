build:
	go build -o ./build/blockx

run: build
	./build/blockx

test:
	go test ./...

clean:
	rm -rf ./build