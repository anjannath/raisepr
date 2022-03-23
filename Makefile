all: raisepr

raisepr: main.go helpers.go types.go
	go build -o $@ $?

.PHONY: clean
clean:
	rm -rf raisepr

.PHONY: test
test:
	@go test -v .