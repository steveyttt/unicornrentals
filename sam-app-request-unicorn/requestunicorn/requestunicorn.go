//https://github.com/ianfoo/aws-lambda-go-demo/blob/master/main.go
//https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-using-start-api.html
//https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-as-simple-proxy-for-lambda.html#api-gateway-create-api-as-simple-proxy-for-lambda-test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type unicorn struct {
	Name   string
	Color  string
	Gender string
}

type unicornRequest struct {
	PickupLocation struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	} `json:"PickupLocation"`
}

//letters to use in the random string generator
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

//(events.APIGatewayProxyResponse, error)

//entry point for lambda
//context always comes first when running a lambda handler
//then the request follows
//context is the
//(events.APIGatewayProxyResponse, error)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) {

	//create unicorns based upon the unicorn struct type
	u1 := unicorn{"Bucephalus", "Golden", "Male"}
	u2 := unicorn{"Shadowfax", "White", "Male"}
	u3 := unicorn{"Rocinante", "Yellow", "Female"}
	//create a variable called unicorns which is a slice of unicorns
	unicorns := []unicorn{u1, u2, u3}
	fmt.Println(unicorns)

	//if the cognito user pool isnt present shit yourself
	if request.RequestContext.Authorizer == nil {
		fmt.Println("AUTHORIZER NOT PRESENT")
	}

	//make a slice of bytes 24 characters long from the const letterbytes
	numRideID := make([]byte, 24)
	for i := range numRideID {
		numRideID[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	rideID := string(numRideID)
	fmt.Println(string(rideID))
	// fmt.Println("RECEIVED EVENT - RIDE ID IS:", rideID, "EVENT ID IS", request)
	// fmt.Println("ONTEXT IS", ctx)

	//Get the username for the session from cognito
	username := request.RequestContext.Authorizer["claims"] //print this in json2go to see if you can get more?
	fmt.Println(username)

	// fmt.Println("REQUEST BODY IS", request.Body)

	//start parsing the request body
	var req unicornRequest

	//unmarshall the shit json into something go can use
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		fmt.Println("HAD PROBLEMS UNMARSHALING", err)
	}

	fmt.Println(req.PickupLocation)

}

//execute lambda
func main() {
	lambda.Start(handler)
}
