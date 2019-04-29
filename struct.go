package main

type Item struct {
	Quarter string // 분기 2019Q1 : Partition Key
	DepositDate string // 예정일 2019-04-11T18:26:00+09:00  : SortKey
	DepositAmount int64 // 입금금액
	ActualDepositDate string // 실입금일
	ActualDepositAmount int64 // 실입금금액
	Typ string // 종류 donation(기부), investment(콘텐츠투자), profit(수익), 계약금(contract),중도금(interim), 잔금(balance), 추가금(addon)
	MonetaryUnit string // 단위 : $, ￦
	Sender string // 보내는이
	Recipient string // 받는이
	Description string // 설명
}
