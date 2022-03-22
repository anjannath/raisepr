all: raisepr

raisepr: main.go
	go build -o $@ $<