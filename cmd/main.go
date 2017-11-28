package main

import (
	"flag"
	// "time"
)

var (
	startDate, endDate, id string
	ErrorInvalidInput      = `Invalid input parameters… please recheck your input`
	ErrorInvalidRange      = `Date range invalid, dates must not be equal`
)

func init() {
	// we use golang flag to parse options
	flag.StringVar(&startDate, `start`, ``, `Start of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&endDate, `end`, ``, `End of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&id, `id`, ``, `User id to fetch invoices`)
	flag.Parse()
}

func main() {

	if startDate == `` || endDate == `` || id == `` {
		panic(ErrorInvalidInput)
	}

	if startDate == endDate {
		panic(ErrorInvalidRange)
	}
}
