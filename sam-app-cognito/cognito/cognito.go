// take imported values from CF as event properties
// replace values inside xs
//push to s3
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event cfn.Event) {
	xs := `
	window._config = {
		cognito: {
			userPoolId: 'replacemeuserPoolId', // e.g. us-east-2_uXboG5pAb
			userPoolClientId: 'replacemeuserPoolClientId', // e.g. 25ddkmj4v6hfsfvruhpfi7n4hv
			region: 'ap-southeast-2' // e.g. us-east-2
		},
		api: {
			invokeUrl: '' // e.g. https://rc7nyt4tql.execute-api.us-west-2.amazonaws.com/prod',
		}
	};
	`

	xs1 := strings.Replace(xs, "replacemeuserPoolId", "STEVEISHERE", 1)
	xs1 = strings.Replace(xs1, "replacemeuserPoolClientId", "BELIEVEINME", 1)

	xs1nr := strings.NewReader(xs1)
	fmt.Printf("%T\n", xs1nr)
	fmt.Println(xs1nr)

	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(xs1)),
		Bucket: aws.String("wildrydes-849375858678"),
		Key:    aws.String("sausage.js"),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		log.Print("SHIT THE BED")
	}
	fmt.Println(result)

}

func main() {
	lambda.Start(handler)
}
