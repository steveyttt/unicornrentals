package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context, event cfn.Event) {

	userPool := event.ResourceProperties["UserPool"]
	us := userPool.(string)
	client := event.ResourceProperties["Client"]
	cs := client.(string)
	bucket := event.ResourceProperties["Bucket"]
	bs := bucket.(string)
	region := event.ResourceProperties["Region"]
	rs := region.(string)

	xs := `
	window._config = {
		cognito: {
			userPoolId: 'replacemeuserPoolId', // e.g. us-east-2_uXboG5pAb
			userPoolClientId: 'replacemeuserPoolClientId', // e.g. 25ddkmj4v6hfsfvruhpfi7n4hv
			region: 'replacemeregion' // e.g. us-east-2
		},
		api: {
			invokeUrl: '' // e.g. https://rc7nyt4tql.execute-api.us-west-2.amazonaws.com/prod',
		}
	};
	`

	xs1 := strings.Replace(xs, "replacemeuserPoolId", us, 1)
	xs1 = strings.Replace(xs1, "replacemeuserPoolClientId", cs, 1)
	xs1 = strings.Replace(xs1, "replacemeregion", rs, 1)

	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(xs1)),
		Bucket: aws.String(bs),
		Key:    aws.String("js/config.js"),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		log.Print(err)
	} else {
		fmt.Println(result)
	}

	r := cfn.NewResponse(&event)
	r.Status = "SUCCESS"
	if r.PhysicalResourceID == "" {
		r.PhysicalResourceID = (lambdacontext.LogStreamName)
	}

	er := r.Send()
	if er != nil {
		fmt.Println(er)
	}

}

func main() {
	lambda.Start(handler)
}
