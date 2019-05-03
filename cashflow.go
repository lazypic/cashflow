package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ladydascalie/currency"
	"github.com/olekukonko/tablewriter"
)

func main() {
	// 인수처리부
	now := time.Now()
	regionPtr := flag.String("region", "ap-northeast-2", "aws region name")
	profilePtr := flag.String("profile", "lazypic", "aws credentials profile name")
	tablePtr := flag.String("table", "cashflow_demo", "aws dynamodb table name")
	datePtr := flag.String("date", now.Format(time.RFC3339), "deposit date")
	amountPtr := flag.Int64("amount", 0, "deposit amount (Required)")
	recipientPtr := flag.String("recipient", "lazypic", "recipient")
	projectPtr := flag.String("project", "none", "project name")
	descriptionPtr := flag.String("description", "none", "description")
	unitPtr := flag.String("unit", "KRW", "mometary unit")
	senderPtr := flag.String("sender", "", "sender (Required)")
	typePtr := flag.String("type", "donation", "type name: donation, investment, profit(일시수익), contract(계약금), interim(중도금), balance(잔금), addon(추가금)")
	actualDatePtr := flag.String("actualdate", now.Format(time.RFC3339), "actual deposit date")
	actualAmountPtr := flag.Int64("actualamount", 0, "actual deposit amount")
	helpPtr := flag.Bool("help", false, "print help")
	flag.Parse()
	if !currency.Valid(*unitPtr) {
		fmt.Fprintf(os.Stderr, "%s string is not ISO4217 format", *unitPtr)
		os.Exit(1)
	}
	if !rfc3339.MatchString(*datePtr) {
		fmt.Fprintf(os.Stderr, "%s string is not RFC3339 format", *datePtr)
		os.Exit(1)
	}
	if !rfc3339.MatchString(*actualDatePtr) {
		fmt.Fprintf(os.Stderr, "%s string is not RFC3339 format", *actualDatePtr)
		os.Exit(1)
	}
	if *helpPtr {
		flag.PrintDefaults()
		os.Exit(0)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*regionPtr)},
		Profile:           *profilePtr,
	}))
	db := dynamodb.New(sess)

	if *senderPtr == "" || *amountPtr == 0 {
		// 분기출력
		for y := now.Year() - 1; y <= now.Year()+2; y++ {
			qr := QuarterlyReport{}
			qr.Year = y
			var totalIn int64
			var totalOut int64
			for q := 1; q <= 4; q++ {
				partitionKey := fmt.Sprintf("%dQ%d", y, q)
				// 분기별 데이터 가지고 와서 출력하기
				in, out, err := QuarterInfo(*db, *tablePtr, partitionKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				switch q {
				case 1:
					qr.Q1.Name = partitionKey
					qr.Q1.In = in
					qr.Q1.Out = out
				case 2:
					qr.Q2.Name = partitionKey
					qr.Q2.In = in
					qr.Q2.Out = out
				case 3:
					qr.Q3.Name = partitionKey
					qr.Q3.In = in
					qr.Q3.Out = out
				case 4:
					qr.Q4.Name = partitionKey
					qr.Q4.In = in
					qr.Q4.Out = out
				}
				totalIn += in
				totalOut += out
			}
			qr.QT.Name = fmt.Sprintf("%dQT", y)
			qr.QT.In = totalIn
			qr.QT.Out = totalOut
			// 테이블 그리기
			data := [][]string{
				[]string{"in",
					strconv.FormatInt(qr.Q1.In, 10),
					strconv.FormatInt(qr.Q2.In, 10),
					strconv.FormatInt(qr.Q3.In, 10),
					strconv.FormatInt(qr.Q4.In, 10),
					strconv.FormatInt(qr.QT.In, 10)},
				[]string{"out",
					strconv.FormatInt(qr.Q1.Out, 10),
					strconv.FormatInt(qr.Q2.Out, 10),
					strconv.FormatInt(qr.Q3.Out, 10),
					strconv.FormatInt(qr.Q4.Out, 10),
					strconv.FormatInt(qr.QT.Out, 10)},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"",
				qr.Q1.Name,
				qr.Q2.Name,
				qr.Q3.Name,
				qr.Q4.Name,
				qr.QT.Name},
			)
			for _, v := range data {
				table.Append(v)
			}
			table.Render() // Send output
		}

		os.Exit(0)
	}

	// 테이블이 존재하는지 점검하고 없다면 테이블을 생성한다.
	if !validTable(*db, *tablePtr) {
		_, err := db.CreateTable(tableStruct(*tablePtr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created the table", *tablePtr)
		fmt.Println("Please re-enter the data after one minute.")
		os.Exit(0)
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
	_, err = db.PutItem(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("add item")
	item.Print()
}
