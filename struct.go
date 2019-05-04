package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Item 은 cashflow에서 사용되는 자료구조이다.
type Item struct {
	Quarter             string  // 분기 2019Q1 : Partition Key
	DepositDate         string  // 예정일 2019-04-11T18:26:00+09:00  : SortKey
	DepositAmount       float64 // 입금금액
	ActualDepositDate   string  // 실입금일
	ActualDepositAmount float64 // 실입금금액
	Typ                 string  // 종류 donation(기부), investment(콘텐츠투자), profit(일시적수익), 계약금(contract),중도금(interim), 잔금(balance), 추가금(addon)
	MonetaryUnit        string  // 단위 : KRW,USD,CNY,JPY,VND / policy : ISO4217
	Sender              string  // 보내는이
	Recipient           string  // 받는이
	Project             string  // 관련 프로젝트명
	Description         string  // 설명
}

// Quarter 는 분기하나의 자료구조이다.
type Quarter struct {
	Name string
	In   float64
	Out  float64
}

// QuarterlyReport 는 분기별 연간 자료구조이다.
type QuarterlyReport struct {
	Year int
	Q1   Quarter
	Q2   Quarter
	Q3   Quarter
	Q4   Quarter
	QT   Quarter
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

// Print 메소드는 QuarterlyReport 자료구조를 표로 출력한다.
func (qr *QuarterlyReport) Print() {
	divValue := 1000000.0
	data := [][]string{
		[]string{"in",
			fmt.Sprintf("%.1f", qr.Q1.In/divValue),
			fmt.Sprintf("%.1f", qr.Q2.In/divValue),
			fmt.Sprintf("%.1f", qr.Q3.In/divValue),
			fmt.Sprintf("%.1f", qr.Q4.In/divValue),
			fmt.Sprintf("%.1f", qr.QT.In/divValue)},
		[]string{"out",
			fmt.Sprintf("%.1f", qr.Q1.Out/divValue),
			fmt.Sprintf("%.1f", qr.Q2.Out/divValue),
			fmt.Sprintf("%.1f", qr.Q3.Out/divValue),
			fmt.Sprintf("%.1f", qr.Q4.Out/divValue),
			fmt.Sprintf("%.1f", qr.QT.Out/divValue)},
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
