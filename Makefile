APP_NAME := kx

.PHONY: build run clean test

build:
	go build -o ./bin/$(APP_NAME) ./main.go

run: build 
	./bin/$(APP_NAME) --shell $$SHELL

clean:
	go clean
	rm -f ./bin/$(APP_NAME)

test:
	go test ./...