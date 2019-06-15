package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/aws/aws-lambda-go/cfn"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-lambda-go/lambda"
)

//copyObjects is designed to copy data from one s3 bucket to another
//it performs some opinionated changes to the folder structure between source and destination
//should take 3 string params
func copyObjects(sb string, sp string, bk string) {

	//start a session to s3
	svc := s3.New(session.New())

	//Create the inputs for api call ListObjectsV2
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(sb),
		Prefix: aws.String(sp),
	}

	//run ListObjectsV2 with the provided above ListObjectsV2Input
	fmt.Println("LISTING OBJECTS")
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		fmt.Println("SOMETHING WENT WRONG PERFORMING ListObjectsV2")
		fmt.Println(err)
	}

	//range over the objects found and copy them to the destination bucket
	for _, v := range result.Contents {

		key := *v.Key                                                                    //key is a reference to the in-memory key value
		source := sb + "/" + key                                                         //source is a variable joining the bucket and key to create one full path to copy from
		copyKey := strings.TrimPrefix(key, "WebApplication/1_StaticWebHosting/website/") //copy key is the destination copy path you want to use - i used trimprefix to strip the path i dont need

		// provide input for the CopyObject call
		input := &s3.CopyObjectInput{
			Bucket:     aws.String(bk),      //my personal bucket
			CopySource: aws.String(source),  //bucket containing files to be copied
			Key:        aws.String(copyKey), //Key - is the location to copy the files
		}

		//call CopyObject with the CopyObjectInput
		_, err := svc.CopyObject(input)
		if err != nil {
			fmt.Println("SOMETHING WENT REALY WRONGWITH THE FILE COPY")
			fmt.Println(err)
		}
	}
}

func handler(ctx context.Context, event cfn.Event) {

	fmt.Println("STARTING LAMBDA --- PRINTING EVENT RECEIVED")
	fmt.Println(event)

	sourceBucket := event.ResourceProperties["SourceBucket"]
	sb := sourceBucket.(string)
	sourcePrefix := event.ResourceProperties["SourcePrefix"]
	sp := sourcePrefix.(string)
	bucket := event.ResourceProperties["Bucket"]
	bk := bucket.(string)

	fmt.Println("CALLING COPYOBJECTS FUNCTION")
	copyObjects(sb, sp, bk)

	//func NewResponse(r *Event) *Response
	r := cfn.NewResponse(&event)

	//set the status to success for the call back to CF
	r.Status = "SUCCESS"

	//CF will shit itself if physical resource ID is empty in the response
	//to get around this I set it to the log stream name
	if r.PhysicalResourceID == "" {
		r.PhysicalResourceID = (lambdacontext.LogStreamName)
	}

	//run the send method on the response
	er := r.Send()
	if er != nil {
		fmt.Println(er)
	}

}

func main() {
	lambda.Start(handler)
}
