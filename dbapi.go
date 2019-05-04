package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func tableStruct(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Quarter"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("DepositDate"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Quarter"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("DepositDate"),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest), // ondemand
		TableName:   aws.String(tableName),
	}
}

func validTable(db dynamodb.DynamoDB, tableName string) bool {
	input := &dynamodb.ListTablesInput{}
	isTableName := false
	// 한번에 최대 100개의 테이블만 가지고 올 수 있다.
	// 한 리전에 최대 256개의 테이블이 존재할 수 있다.
	// https://docs.aws.amazon.com/ko_kr/amazondynamodb/latest/developerguide/Limits.html
	for {
		result, err := db.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Fprintf(os.Stderr, "%s %s\n", dynamodb.ErrCodeInternalServerError, err.Error())
				default:
					fmt.Fprintf(os.Stderr, "%s\n", aerr.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			return false
		}

		for _, n := range result.TableNames {
			if *n == tableName {
				isTableName = true
				break
			}
		}
		if isTableName {
			break
		}
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
	return isTableName
}

// GetQuarter 함수는 "2019Q1" 형태의 문자를 입력받아서 수입,지출 정보를 가지고 온다.
func GetQuarter(db dynamodb.DynamoDB, tableName string, quarter string) (int64, int64, error) {
	var in int64
	var out int64

	filt := expression.Name("Quarter").Equal(expression.Value(quarter))
	proj := expression.NamesList(expression.Name("DepositAmount"), expression.Name("Sender"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return in, out, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := db.Scan(params)
	if err != nil {
		return in, out, err
	}

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return in, out, err
		}
		if item.Sender == "lazypic" {
			out += item.DepositAmount
		}
		in += item.DepositAmount
	}
	return in, out, err
}

func hasItem(db dynamodb.DynamoDB, tableName string, primarykey string, sortkey string) bool {
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Quarter": {
				S: aws.String(primarykey),
			},
			"DepositDate": {
				S: aws.String(sortkey),
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return false
	}
	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return false
	}
	if item.DepositDate == "" {
		return false
	}
	return true
}
