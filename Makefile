build:
	@go build -o bin/ecommerce_backend
	
run:	build
	@./bin/ecommerce_backend

test:
	@go test -v ./...   
