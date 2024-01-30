PHONY: build-cli
build-cli:
	go build 


PHONY: install-cli
install-cli: build-cli
	sudo ln -sf ${PWD}/hc /usr/local/bin/hc


PHONY: build-image
build-image: build-cli
	./hc build


PHONY: all
all: build-cli install-cli build-image
