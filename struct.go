package main

type Item struct {
	Quarter string // 분기 2019Q1 : Partition Key
	DepositDate string // 예정일 2019-04-11T18:26:00+09:00  : SortKey
	DepositAmount int64 // 입금금액
	ActualDepositDate string // 실입금일
	ActualDepositAmount int64 // 실입금금액
	Typ string // 종류 donation, investment, profit
	MonetaryUnit string // 단위 : $, ￦
	Sender string // 보내는이
	Recipient string // 받는이
	Description string // 설명
}
