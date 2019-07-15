.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./request-unicorn/requestunicorn
	
build:
	GOOS=linux GOARCH=amd64 go build -o requestunicorn/requestunicorn ./requestunicorn

run:
	sam local invoke requestunicorn -e sampleevent.json

package:
	sam package --output-template-file packaged.yaml --s3-bucket lambda-849375858678 --template-file template.yaml

deploy:
	sam deploy --template-file ./packaged.yaml --stack-name request-Unicorn --capabilities CAPABILITY_IAM