build:
	go build -o ~/bin/i3-bar main.go

try: build
	i3-msg restart

dev:
	docker run --init -it --rm --network none \
	-v $(shell pwd):/i3-bar \
	golang:1.12.3-stretch /bin/bash
