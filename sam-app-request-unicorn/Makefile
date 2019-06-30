.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./request-unicorn/request-unicorn
	
build:
	GOOS=linux GOARCH=amd64 go build -o request-unicorn/request-unicorn ./request-unicorn

run:
	sam local invoke request-unicorn -e sampleevent.json