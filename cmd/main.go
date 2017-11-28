package main

import (
	"flag"
	"fmt"
	"sync"
	// "time"
)

const (
	ErrorInvalidInput = `Invalid input parameters… please recheck your input`
	ErrorInvalidRange = `Date range invalid, dates must not be equal`
	// API base url
	baseURL       = `http://34.209.24.195/facturas`
	requestFormat = baseURL + `?id=%s&start=%s&end=%s`
)

var (
	startDate, endDate, id     string
	requestCount, invoiceCount int
	// I'll use a waitgroup to wait for requests to be processed
	wg sync.WaitGroup
)

func init() {
	// I use golang flag to parse command line options
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

	fetchInvoices(id, startDate, endDate)

	wg.Wait()
	fmt.Printf("%d invoices were found, using %d requests\n", invoiceCount, requestCount)
}

func fetchInvoices(id, start, end string) {
	// Increment waitgroup count
	wg.Add(1)
	requestCount += 1
	// TBD
}
