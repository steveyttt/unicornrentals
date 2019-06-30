package main

import (
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

//letters to use in the ranfom string generator
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

//entry point for lambda
func handler(request events.APIGatewayProxyRequest) {

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

	fmt.Println(request.RequestContext.Authorizer)

	//make a slice of bytes 24 characters long from the const letterbytes
	b := make([]byte, 24)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	fmt.Println(string(b))

}

//execute lambda
func main() {
	lambda.Start(handler)
}
