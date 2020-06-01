package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type server struct{}

// RequestData struct to parse incoming json
type RequestData struct {
	Type      string
	DateStart int
	DateEnd   int
}

// ResponseData struct to parse outcoming json
type ResponseData struct {
	Id    string
	Date  int
	Value float64
	Type  string
}

func main() {
	//define routes
	http.HandleFunc("/getRange", getRange)
	http.HandleFunc("/getLast", getLast)
	http.HandleFunc("/healthCheck", healthCheck)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//parse request parse the incoming json return the serialized object back
func parseRequest(r *http.Request) (rd RequestData, e error) {
	var req RequestData

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, errors.New(err.Error())
	}

	return req, nil
}

//health check return web services status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": true}`))
}

//get the most recent inserted data
func getLast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqData, e := parseRequest(r)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := "Sprinkler"

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"Type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(reqData.Type),
					},
				},
			},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int64(1),
	}

	result, err := svc.Query(queryInput)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return
	}

	item := ResponseData{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	b, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//get a set of data based on given request json
func getRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RequestData

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	svc := dynamodb.New(sess)

	tableName := "Sprinkler"

	// Set filter with given range
	filt := expression.Name("Date").Between(expression.Value(req.DateStart), expression.Value(req.DateEnd))

	proj := expression.NamesList(expression.Name("Date"), expression.Name("Id"), expression.Name("Value"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
	}

	var elements []ResponseData

	for index, element := range result.Items {
		item := ResponseData{}

		err = dynamodbattribute.UnmarshalMap(result.Items[index], &item)
		if err != nil {
			println("error in element %e", element)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		elements = append(elements, item)
	}

	b, err := json.Marshal(elements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
