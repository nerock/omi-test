test:
	go test -v ./...

test-integration:
	go -v -tags integration ./...

build:
	go build -o bin/auditlogger cmd/main.go

build-docker:
	docker build -t auditlogger .