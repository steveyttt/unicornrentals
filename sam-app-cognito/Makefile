.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./cognito/cognito
	
build:
	GOOS=linux GOARCH=amd64 go build -o cognito/cognito ./cognito

run:
	sam local invoke cognito -e sampleevent.json