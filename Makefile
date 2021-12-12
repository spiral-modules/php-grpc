all:
	@./build.sh all
build:
	@./build.sh
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
	rm -rf coverage-ci
	mkdir ./coverage-ci

	go test -v -race -cover -tags=debug -coverpkg=./... -failfast -coverprofile=./coverage-ci/root.out -covermode=atomic .
	go test -v -race -cover -tags=debug -coverpkg=./... -failfast -coverprofile=./coverage-ci/parser.out -covermode=atomic ./parser
	go test -v -race -cover -tags=debug -coverpkg=./... -failfast -coverprofile=./coverage-ci/gen.out -covermode=atomic ./cmd/protoc-gen-php-grpc
	echo 'mode: atomic' > ./coverage-ci/summary.txt
	tail -q -n +2 ./coverage-ci/*.out >> ./coverage-ci/summary.txt

	vendor_php/bin/phpunit

lint:
	go fmt ./...
	golint ./...