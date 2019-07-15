aws cloudformation deploy \
    --stack-name s3-lambda \
    --tags $(cat tags.properties) \
    --template-file ../cf/lambda.s3.yaml

# run me from the deploy folder
aws cloudformation validate-template --template-body file://../cf/s3.yaml
aws cloudformation validate-template --template-body file://../cf/lambda.s3.yaml
aws cloudformation validate-template --template-body file://../cf/cognito.yaml

cd ../sam-app-s3
make build
cd s3copy
zip s3copy.zip s3copy
aws s3 cp s3copy.zip s3://lambda-849375858678/s3copy.zip
cd ../../deploy/

aws cloudformation deploy \
    --stack-name unicorn-rentals-s3 \
    --tags $(cat tags.properties) \
    --parameter-overrides $(cat unicorn-rentals-s3.properties) \
    --template-file ../cf/s3.yaml \
    --capabilities CAPABILITY_NAMED_IAM

cd ../sam-app-cognito
make build
cd cognito
zip cognito.zip cognito
aws s3 cp cognito.zip s3://lambda-849375858678/cognito.zip
cd ../../deploy/

aws cloudformation deploy \
    --stack-name unicorn-rentals-cognito \
    --tags $(cat tags.properties) \
    --parameter-overrides $(cat unicorn-rentals-cognito.properties) \
    --template-file ../cf/cognito.yaml \
    --capabilities CAPABILITY_NAMED_IAM

cd ../sam-app-request-unicorn
make package
make deploy
cd ../deploy
