export CGO_ENABLED ?= 0

.PHONY: build
build:
	go build -o since

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: run
run:
	go run ./main.go $(filter-out $@,$(MAKECMDGOALS))

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
