all:
	@./build.sh
build:
	@./build.sh all
clean:
	rm -rf protoc-gen-php-grpc
	rm -rf rr-grpc
install: all
	cp rr /usr/local/bin/protoc-gen-php-grpc
	cp rr /usr/local/bin/rr-grpc
uninstall: 
	rm -f /usr/local/bin/protoc-gen-php-grpc
	rm -f /usr/local/bin/rr-grpc
test:
	go test -v -race -cover
	go test -v -race -cover ./parser