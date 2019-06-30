.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./request-unicorn/requestunicorn
	
build:
	GOOS=linux GOARCH=amd64 go build -o requestunicorn/requestunicorn ./requestunicorn

run:
	sam local invoke requestunicorn -e sampleevent.json