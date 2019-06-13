.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./s3copy/s3copy
	
build:
	GOOS=linux GOARCH=amd64 go build -o s3copy/s3copy ./s3copy

run:
	sam local invoke s3copy -e sampleevent.json