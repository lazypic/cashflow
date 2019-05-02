package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ladydascalie/currency"
)

func main() {
	// 인수처리부
	now := time.Now().Format(time.RFC3339)
	regionPtr := flag.String("region", "ap-northeast-2", "aws region name")
	profilePtr := flag.String("profile", "lazypic", "aws credentials profile name")
	tablePtr := flag.String("table", "cashflow_demo", "aws dynamodb table name")
	datePtr := flag.String("date", now, "deposit date")
	amountPtr := flag.Int64("amount", 0, "deposit amount (Required)")
	recipientPtr := flag.String("recipient", "lazypic", "recipient")
	projectPtr := flag.String("project", "none", "project name")
	descriptionPtr := flag.String("description", "none", "description")
	unitPtr := flag.String("unit", "KRW", "mometary unit")
	senderPtr := flag.String("sender", "", "sender (Required)")
	typePtr := flag.String("type", "donation", "type name: donation, investment, profit(일시수익), contract(계약금), interim(중도금), balance(잔금), addon(추가금)")
	actualDatePtr := flag.String("actualdate", now, "actual deposit date")
	actualAmountPtr := flag.Int64("actualamount", 0, "actual deposit amount")
	helpPtr := flag.Bool("help", false, "print help")
	flag.Parse()
	if !currency.Valid(*unitPtr) {
		fmt.Fprintf(os.Stderr, "%s string is not ISO4217 format", *unitPtr)
		os.Exit(0)
	}
	if *helpPtr {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *senderPtr == "" || *amountPtr == 0 {
		// 분기별 데이터 가지고 오기
		// 분기별 보고 출력하기.
		fmt.Println("drawing chart")
		os.Exit(0)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*regionPtr)},
		Profile:           *profilePtr,
	}))
	svc := dynamodb.New(sess)

	// 테이블이 존재하는지 점검하고 없다면 테이블을 생성한다.
	if !validTable(*svc, *tablePtr) {
		_, err := svc.CreateTable(tableStruct(*tablePtr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created the table", *tablePtr)
	}
	// 아이템 추가하기.
	q, err := RFC3339ToQuarter(*datePtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	if *actualAmountPtr == 0 {
		*actualAmountPtr = *amountPtr
	}

	item := Item{
		Quarter:             q,
		DepositDate:         *datePtr,
		DepositAmount:       *amountPtr,
		ActualDepositDate:   *actualDatePtr,
		ActualDepositAmount: *actualAmountPtr,
		Typ:                 *typePtr,
		MonetaryUnit:        *unitPtr,
		Sender:              *senderPtr,
		Recipient:           *recipientPtr,
		Project:             *projectPtr,
		Description:         *descriptionPtr,
	}
	dynamodbJSON, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*tablePtr),
	}
	_, err = svc.PutItem(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("add item")
	item.Print()
}
