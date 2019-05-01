package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "cashflow_demo"

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("ap-northeast-2")},
		Profile:           "lazypic",
	}))
	svc := dynamodb.New(sess)

	// 테이블 리스트 체크하기.
	if !isTableName(*svc, tableName) {
		_, err := svc.CreateTable(tableStruct(tableName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created the table", tableName)
	}
	// argv받기.
	// 아이템 추가하기.
	item := Item{
		Quarter:             "2019Q2",
		DepositDate:         "2019-06-12T18:26:00+09:00",
		DepositAmount:       10000,
		ActualDepositDate:   "2019-06-12T18:26:00+09:00",
		ActualDepositAmount: 10000,
		Typ:                 "donation",
		MonetaryUnit:        "KRW",
		Sender:              "test",
		Recipient:           "lazypic",
		Project:             "project name",
		Description:         "description",
	}
	dynamodbJSON, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully added %v\n", item)
	// 데이터 가지고 오기
	// 분기별 보고 출력하기.
}
