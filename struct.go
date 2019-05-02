package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Item 은 cashflow에서 사용되는 자료구조이다.
type Item struct {
	Quarter             string // 분기 2019Q1 : Partition Key
	DepositDate         string // 예정일 2019-04-11T18:26:00+09:00  : SortKey
	DepositAmount       int64  // 입금금액
	ActualDepositDate   string // 실입금일
	ActualDepositAmount int64  // 실입금금액
	Typ                 string // 종류 donation(기부), investment(콘텐츠투자), profit(일시적수익), 계약금(contract),중도금(interim), 잔금(balance), 추가금(addon)
	MonetaryUnit        string // 단위 : KRW,USD,CNY,JPY,VND / policy : ISO4217
	Sender              string // 보내는이
	Recipient           string // 받는이
	Project             string // 관련 프로젝트명
	Description         string // 설명
}

// Print 메소드는 Item 자료구조를 보기좋게 출력한다.
func (i *Item) Print() {
	fmt.Println("Quarter:", i.Quarter)
	fmt.Println("Deposit Date:", i.DepositDate)
	fmt.Println("Deposit Amount:", i.MonetaryUnit, i.DepositAmount)
	fmt.Println("Actual Deposit Date:", i.ActualDepositDate)
	fmt.Println("Actual Deposit Amount:", i.MonetaryUnit, i.ActualDepositAmount)
	fmt.Println("Type:", i.Typ)
	fmt.Printf("%s -> %s\n", i.Sender, i.Recipient)
	fmt.Println("Project:", i.Project)
	fmt.Println("Description:", i.Description)
}

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
