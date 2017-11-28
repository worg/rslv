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

	// Process the increment of invoices & request count on channels
	requestChan := make(chan int)
	invoiceChan := make(chan int)

	// Wait for increment values on a goroutine
	go processIncrements(requestChan, invoiceChan)

	fetchInvoices(id, startDate, endDate, requestChan, invoiceChan)

	wg.Wait()
	close(requestChan)
	close(invoiceChan)

	fmt.Printf("%d invoices were found, using %d requests\n", invoiceCount, requestCount)
}

func fetchInvoices(id, start, end string, requestChan, invoiceChan chan int) {
	// Increment waitgroup count
	wg.Add(1)
	requestChan <- 1
	defer wg.Done()
	// TBD
}

// Read increments from channels
func processIncrements(requestChan, invoiceChan chan int) {
	select {
	case <-requestChan:
		requestCount += 1
	case add, ok := <-invoiceChan:
		if !ok {
			break
		}
		invoiceCount += add
	}
}
