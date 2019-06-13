# run me from the deploy folder
aws cloudformation validate-template --template-body file://../cf/s3.yaml
aws cloudformation validate-template --template-body file://../cf/lambda.s3.yaml
aws cloudformation validate-template --template-body file://../cf/cognito.yaml

aws cloudformation deploy \
    --stack-name s3-lambda \
    --tags $(cat tags.properties) \
    --template-file ../cf/lambda.s3.yaml

cd ../sam-app-s3
make build
cd s3copy
zip s3copy.zip s3copy
aws s3 cp s3copy.zip s3://lambda-849375858678/s3copy.zip
cd ../../deploy/

aws cloudformation deploy \
    --stack-name s3-unicorn-rentals \
    --tags $(cat tags.properties) \
    --parameter-overrides $(cat s3-unicorn-rentals-s3.properties) \
    --template-file ../cf/s3.yaml \
    --capabilities CAPABILITY_NAMED_IAM

aws cloudformation deploy \
    --stack-name s3-unicorn-cognito \
    --tags $(cat tags.properties) \
    --parameter-overrides $(cat s3-unicorn-rentals-cognito.properties) \
    --template-file ../cf/cognito.yaml \
    --capabilities CAPABILITY_NAMED_IAM
