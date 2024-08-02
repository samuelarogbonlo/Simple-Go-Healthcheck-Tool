EXEC = health-check-tool

SRC = main.go

all: run

build:
	go build -o $(EXEC) $(SRC)

run: build
	./$(EXEC)

clean:
	rm -f $(EXEC) report.json

.PHONY: all build run clean
