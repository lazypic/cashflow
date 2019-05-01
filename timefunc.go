package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var rfc3339 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$`)
var quarter = regexp.MustCompile(`^\d{4}Q\d{1}$`)
var shortQuarter = regexp.MustCompile(`^\d{1}Q\d{2}$`)

// RFC3339ToQuarter 함수는 RFC3339 시간형식을 쿼터표기로 변환한다.
func RFC3339ToQuarter(rfctime string) (string, error) {
	if !rfc3339.MatchString(rfctime) {
		return rfctime, errors.New("time string is not rfc3339 format")
	}
	t, err := time.Parse(time.RFC3339, rfctime)
	if err != nil {
		return rfctime, err
	}
	return fmt.Sprintf("%dQ%d", t.Year(), ((t.Month()-1)/3)+1), nil
}

// RFC3339ToShortQuarter 함수는 RFC3339 시간형식을 짧은 쿼터표기로 변환한다.
func RFC3339ToShortQuarter(rfctime string) (string, error) {
	if !rfc3339.MatchString(rfctime) {
		return rfctime, errors.New("time string is not RFC3339 format")
	}
	t, err := time.Parse(time.RFC3339, rfctime)
	if err != nil {
		return rfctime, err
	}
	return fmt.Sprintf("%dQ%s", ((t.Month()-1)/3)+1, rfctime[2:4]), nil
}
