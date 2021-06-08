.PHONY: mocks

lint:
	golangci-lint run

test:
	go test ./...

format:
	go fmt ./...

update:
	go get -u -d ./... && go mod tidy
