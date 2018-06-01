package db

import (
	_ "fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DbHelper struct {
	Region string
	client *dynamodb.DynamoDB
}

func New(region string) *DbHelper {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return &DbHelper{
		Region: region,
		client: svc,
	}
}
