package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	// ErrorInvalidInput is returned when any required parameter is empty
	ErrorInvalidInput = `Invalid input parameters… please recheck your input`
	// ErrorInvalidRange is returned when trying to get the same date
	ErrorInvalidRange = `Date range invalid, dates must not be equal`
	// ErrorInvalidFormat is returned when dates did not conform the specified format
	ErrorInvalidFormat = `Date format must be YYYY-MM-DD`
	// ErrorInvalidStart is returned when start date is greater than end in range
	ErrorInvalidStart = `Date range invalid, start must not be greater than end`
	// API base url
	baseURL       = `http://34.209.24.195/facturas`
	requestFormat = baseURL + `?id=%s&start=%s&end=%s`
	// Time format allowed
	timeFmt = `2006-01-02`
)

var (
	startDate, endDate, id     string
	requestCount, invoiceCount int
)

func init() {
	// I use golang flag to parse command line options
	flag.StringVar(&startDate, `start`, ``, `Start of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&endDate, `end`, ``, `End of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&id, `id`, ``, `User id to fetch invoices`)
	flag.Parse()
}

func main() {
	var start, end time.Time
	var err error

	if startDate == `` || endDate == `` || id == `` {
		panic(ErrorInvalidInput)
	}

	if startDate == endDate {
		panic(ErrorInvalidRange)
	}

	if start, err = time.Parse(timeFmt, startDate); err != nil {
		panic(ErrorInvalidFormat)
	}

	if end, err = time.Parse(timeFmt, endDate); err != nil {
		panic(ErrorInvalidFormat)
	}

	// Validate that start is < than end
	if start.After(end) {
		panic(ErrorInvalidStart)
	}

	fmt.Printf("%d invoices were found, using %d requests\n", invoiceCount, requestCount)
}
