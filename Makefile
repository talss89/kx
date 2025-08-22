APP_NAME := kx

.PHONY: build run clean test lint build-cov integration

build:
	go build -o ./bin/$(APP_NAME) ./main.go

build-cov:
	go build -o ./bin/$(APP_NAME)-cov -cover ./main.go

run: build
	./bin/$(APP_NAME) --shell $$SHELL

clean:
	go clean
	rm -f ./bin/$(APP_NAME)

unit:
	go test ./internal/... -cover

test: build-cov
	@rm -rf .coverdata
	@mkdir -p .coverdata/integration
	@mkdir -p .coverdata/unit
	@mkdir -p .coverdata/merged
	go test ./... -cover -args -test.gocoverdir=${PWD}/.coverdata/unit
	@go tool covdata merge -i=.coverdata/unit,.coverdata/integration -o=.coverdata/merged
	@go tool covdata percent -i=.coverdata/merged
	@go tool covdata textfmt -i=.coverdata/merged -o .coverdata/coverage.txt
	go tool cover -func=.coverdata/coverage.txt

lint:
	golangci-lint run ./...
