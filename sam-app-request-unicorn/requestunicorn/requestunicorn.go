//https://github.com/ianfoo/aws-lambda-go-demo/blob/master/main.go
//https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-using-start-api.html
//https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-as-simple-proxy-for-lambda.html#api-gateway-create-api-as-simple-proxy-for-lambda-test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Create a struct to store the unicorn informatiom
type unicorn struct {
	Name   string
	Color  string
	Gender string
}

//A struct to hold the infomation relating to the request of a unicorn
type unicornRequest struct {
	PickupLocation struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	} `json:"PickupLocation"`
}

//A struct to hold the responses back to the API gateway
type unicornResponse struct {
	RideID  string
	Rider   interface{}
	Unicorn unicorn
	ETA     string
}

//A struct to hold the data type to inject into dynamoDB
type unicornRecord struct {
	RideID      string
	User        interface{}
	Unicorn     unicorn
	RequestTime time.Time
}

//letters to use in the random string generator
const letterBytes = "abcdefghij1klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

//entry point for lambda
//context always comes first when running a lambda handler
//then the request follows
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//create unicorns based upon the unicorn struct type
	u1 := unicorn{"Bucephalus", "Golden", "Male"}
	u2 := unicorn{"Shadowfax", "White", "Male"}
	u3 := unicorn{"Rocinante", "Yellow", "Female"}
	//create a variable called unicorns which is a slice of unicorns
	//then randomly pick one for use
	rand.Seed(time.Now().UnixNano()) //https://gobyexample.com/random-numbers
	unicorns := []unicorn{u1, u2, u3}
	randint := rand.Intn(3)
	myUnicorn := unicorns[randint]

	//determine time and store as a variable to add to the DB when logging the request
	time := time.Now()

	//make a slice of bytes 24 characters long from the const letterbytes
	//This generates a random string which we can use as a rideID when adding the request to the database
	numRideID := make([]byte, 24)
	for i := range numRideID {
		numRideID[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	rideID := string(numRideID)
	fmt.Println(string(rideID))

	//Get the username for the session from cognito
	claim := request.RequestContext.Authorizer["claims"]
	username := claim.(map[string]interface{})["cognito:username"]

	//start a session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//start a connector to the dynamo api using the above session
	svc := dynamodb.New(sess)

	//create a record for the put item api call to dynamo
	recordrequest := unicornRecord{
		RideID:      rideID,
		User:        username,
		Unicorn:     myUnicorn,
		RequestTime: time,
	}

	//define the dynamo table name
	tableName := "Rides"

	//marshall the values into something go and dynamo can understand
	av, err := dynamodbattribute.MarshalMap(recordrequest)
	if err != nil {
		fmt.Println("Got error marshalling item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//call the put item command and reference the table name + the marshalled item to add
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//Create a response to send to the api gateway returner
	resp := unicornResponse{
		RideID:  rideID,
		Unicorn: myUnicorn,
		Rider:   username,
		ETA:     "30 seconds",
	}

	//put the response into json to send back to client
	reqBytes, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}
	//send the response to the apigateway
	return events.APIGatewayProxyResponse{
		Body:       string(reqBytes),
		StatusCode: http.StatusOK,
	}, nil

}

//execute lambda
func main() {
	lambda.Start(handler)
}
