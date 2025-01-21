build:
	go build -o bin/go-ecommerce-api

run: build
	./bin/go-ecommerce-api

clean:
	rm -rf bin/go-ecommerce-api

test:
	go test -v ./...