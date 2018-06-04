package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type DbHelper struct {
	Region string
	client *dynamodb.DynamoDB
}

// For singleton
var instance *DbHelper
var once sync.Once

// region is blank or aws-region such as ap-northeast-2
// In case of blank, use AWS_DEFAULT_REGION as region
func GetInstance(region string) *DbHelper {
	once.Do(func() {
		instance = New(region)
	})
	return instance
}

func New(region string) *DbHelper {
	// Check if instance is already assigned
	if instance != nil {
		return instance
	}

	// Create AWS session
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	},
	)
	if err != nil {
		fmt.Println("Failed to create AWS session")
		return nil
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	if svc == nil {
		fmt.Println("Failed to create DynamoDB client")
		return nil
	}

	// Assign
	instance = &DbHelper{
		Region: region,
		client: svc,
	}
	return instance
}

func (d *DbHelper) ListTables() {
	result, err := d.client.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Tables:")
	for _, n := range result.TableNames {
		fmt.Println(*n)
	}
}

// For table named "tblName" in DynamoDB,
// Scan target value following propName is propVal
// -----------------------------
// |  propName  |  targetName  |
// -----------------------------
// |  propVal   |  targetVal   |
// |   ...      |    ...       |
// -----------------------------
func (d *DbHelper) GetItem(tblName, propName, propVal, targetName string) *dynamodb.ScanOutput {
	/*
		result, err := d.client.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tblName),
			Key: map[string]*dynamodb.AttributeValue{
				propName: {
					S: aws.String(propVal),
				},
			},
		})
	*/
	filter := expression.Name(propName).Equal(expression.Value(propVal))
	projection := expression.NamesList(expression.Name(propName), expression.Name(targetName))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tblName),
	}

	// Make the DynamoDB Query API call
	result, err := d.client.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return result
}

func (d *DbHelper) UnmarshalMap(in map[string]*dynamodb.AttributeValue, out interface{}) {
	if in == nil {
		return
	} else if err := dynamodbattribute.UnmarshalMap(in, &out); err != nil {
		fmt.Println("Failed to unmarshalMap of dynamoDb output")
	}
}
