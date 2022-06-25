.PHONY: all
all: build

.PHONY: build
build:
	@go build
	@./lts build:done --no-successful-console

.PHONY: clean
clean:
	rm ./lts