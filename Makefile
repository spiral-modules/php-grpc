all:
	@./build.sh
build:
	@./build.sh all
clean:
	rm -rf protoc-gen-php-grpc
	rm -rf rr-grpc
install: all
	cp protoc-gen-php-grpc /usr/local/bin/protoc-gen-php-grpc
	cp rr-grpc /usr/local/bin/rr-grpc
uninstall: 
	rm -f /usr/local/bin/protoc-gen-php-grpc
	rm -f /usr/local/bin/rr-grpc
test:
	go test -v -race -cover
	go test -v -race -cover ./parser
	go test -v -race -cover ./cmd/protoc-gen-php-grpc