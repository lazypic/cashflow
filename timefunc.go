package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var RFC3339 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$`)
var Quarter = regexp.MustCompile(`^\d{4}Q\d{1}$`)
var ShortQuarter = regexp.MustCompile(`^\d{1}Q\d{2}$`)

// RFC3339 시간형식을 분기로 변환한다.₩
func RFC3339_to_Quarter(rfctime string) (string, error) {
	if !RFC3339.MatchString(rfctime) {
		return rfctime, errors.New("time string is not RFC3339 format.")
	}
	t, err := time.Parse(time.RFC3339, rfctime)
	if err != nil {
		return rfctime, err
	}
	return fmt.Sprintf("%dQ%d", t.Year(), t.Month()/3+1), nil
}