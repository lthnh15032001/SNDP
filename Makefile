NAME=github.com/lthnh15032001/ngrok-impl
VERSION=0.0.1
VERSION_PACKAGE=github.com/lthnh15032001/ngrok-impl/internal/version

.PHONY: build
build:
	@echo Building from source....
	@CGO_ENABLED=0 go build -o ./build/$(NAME) ./cmd

.PHONY: run
run: build
	@echo Starting your app using dev configs....
	@./build/$(NAME) -e dev

.PHONY: build-prod
build-prod:
	@echo Building from source....
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

.PHONY: run-prod
run-prod:
	@echo Starting app using prod configs....
	@$CI_PROJECT_DIR/golang-rest-api-starter-binary -e dev

.PHONY: clean
clean:
	@echo Removing build file....
	@rm -f ./build/$(NAME)

.PHONY: test
test:
	@go test -v ./tests/*
