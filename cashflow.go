package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

func main() {
	// 테이블 리스트 체크하기.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("ap-northeast-2")},
		Profile:           "lazypic",
	}))
	svc := dynamodb.New(sess)
	input := &dynamodb.ListTablesInput{}
	tableName := "cashflow_demo"
	isTableName := false
	// 한번에 최대 100개의 테이블만 가지고 올 수 있다.
	// 한 리전에 최대 256개의 테이블이 존재할 수 있다.
	// https://docs.aws.amazon.com/ko_kr/amazondynamodb/latest/developerguide/Limits.html
	for {
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
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

	if !isTableName {
		_, err := svc.CreateTable(tableStruct(tableName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created the table", tableName)
	}
	// argv받기.
	// 구현필요함.
	// 아이템 추가하기.
	item := Item{
		Quarter:             "2019Q1",
		DepositDate:         "2019-04-12T18:26:00+09:00",
		DepositAmount:       10000,
		ActualDepositDate:   "2019-04-12T18:26:00+09:00",
		ActualDepositAmount: 10000,
		Typ:                 "donation",
		MonetaryUnit:        "KRW",
		Sender:              "test",
		Recipient:           "lazypic",
		Project:             "project name",
		Description:         "description",
	}
	dynamodbJson, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJson,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully added %v\n", item)
	// 분기별 보고 출력하기.
}
