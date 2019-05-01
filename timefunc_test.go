package main

import (
	"testing"
)

func Test_RFC3339_to_Quarter(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{
		in:   "2019-01-10T10:00:00+09:00",
		want: "2019Q1",
	}, {
		in:   "2019-02-28T10:00:00+09:00",
		want: "2019Q1",
	}, {
		in:   "2019-03-30T10:00:00+09:00",
		want: "2019Q1",
	}, {
		in:   "2019-04-30T10:00:00+09:00",
		want: "2019Q2",
	}, {
		in:   "2019-05-30T10:00:00+09:00",
		want: "2019Q2",
	}, {
		in:   "2019-06-30T10:00:00+09:00",
		want: "2019Q2",
	}, {
		in:   "2019-07-30T10:00:00+09:00",
		want: "2019Q3",
	}, {
		in:   "2019-08-30T10:00:00+09:00",
		want: "2019Q3",
	}, {
		in:   "2019-09-30T10:00:00+09:00",
		want: "2019Q3",
	}, {
		in:   "2019-10-30T10:00:00+09:00",
		want: "2019Q4",
	}, {
		in:   "2019-11-30T10:00:00+09:00",
		want: "2019Q4",
	}, {
		in:   "2019-12-30T10:00:00+09:00",
		want: "2019Q4",
	}}
	for _, c := range cases {
		got, err := RFC3339_to_Quarter(c.in)
		// check err
		if err != nil {
			t.Fatalf("%v", err.Error())
		}
		// check value
		if got != c.want {
			t.Fatalf("RFC3339_to_Quarter(%v): got %s, want %s", c.in, got, c.want)
		}
	}
}

func Test_RFC3339_to_ShortQuarter(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{
		in:   "2019-01-10T10:00:00+09:00",
		want: "1Q19",
	}, {
		in:   "2019-02-28T10:00:00+09:00",
		want: "1Q19",
	}, {
		in:   "2019-03-30T10:00:00+09:00",
		want: "1Q19",
	}, {
		in:   "2019-04-30T10:00:00+09:00",
		want: "2Q19",
	}, {
		in:   "2019-05-30T10:00:00+09:00",
		want: "2Q19",
	}, {
		in:   "2019-06-30T10:00:00+09:00",
		want: "2Q19",
	}, {
		in:   "2019-07-30T10:00:00+09:00",
		want: "3Q19",
	}, {
		in:   "2019-08-30T10:00:00+09:00",
		want: "3Q19",
	}, {
		in:   "2019-09-30T10:00:00+09:00",
		want: "3Q19",
	}, {
		in:   "2019-10-30T10:00:00+09:00",
		want: "4Q19",
	}, {
		in:   "2019-11-30T10:00:00+09:00",
		want: "4Q19",
	}, {
		in:   "2019-12-30T10:00:00+09:00",
		want: "4Q19",
	}}
	for _, c := range cases {
		got, err := RFC3339_to_ShortQuarter(c.in)
		// check err
		if err != nil {
			t.Fatalf("%v", err.Error())
		}
		// check value
		if got != c.want {
			t.Fatalf("RFC3339_to_ShortQuarter(%v): got %s, want %s", c.in, got, c.want)
		}
	}
}
