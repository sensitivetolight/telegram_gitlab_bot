.SILENT:

SHELL := /bin/bash

PORT=$(shell cat config/port)
BINARY=gitlab_telegram_bot

build:
			go build -o $(BINARY)

run:
			go run telegram.go 

test: 
			for filename in testfiles/*.json; do \
  			curl --data-binary "@$$filename" -H "Content-Type: application/json" -X POST http\://localhost\:22222 ; \
			done 

docker:
			echo "Building docker image with name $name"
			docker build -t $(BINARY) . 
			echo "Running docker image and expose port 22222"
			docker run -d -p $(PORT):22222 $(BINARY)

clean:
			rm $(BINARY)

