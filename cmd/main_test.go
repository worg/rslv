package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestResuelve(t *testing.T) {
	Convey(`When running the app`, t, func() {
		Convey(`Input must not be empty`, func() {
			So(main, ShouldPanicWith, ErrorInvalidInput)
		})

		Convey(`Dates must differ`, func() {
			startDate = `2017-11-01`
			endDate = `2017-11-01`
			id = `1`
			So(main, ShouldPanicWith, ErrorInvalidRange)
		})

		Convey(`Dates must comply with the specified format`, func() {
			startDate = `2017 Jan 01`
			endDate = `2017 Jan 03`
			id = `1`
			So(main, ShouldPanicWith, ErrorInvalidFormat)
		})

		Convey(`When making requests`, func() {
			startDate = `2017-01-01`
			endDate = `2017-03-01`
			id = `1`
			start, _ := time.Parse(timeFmt, startDate)
			end, _ := time.Parse(timeFmt, endDate)
			var requestChan, invoiceChan chan int

			reset := func() {
				requestChan = make(chan int)
				invoiceChan = make(chan int)
				requestCount = 0
				invoiceCount = 0
			}
			reset()

			Convey(`Counters should start on zero`, func() {
				So(requestCount, ShouldEqual, 0)
				So(invoiceCount, ShouldEqual, 0)
			})

			Convey(`Should Increment the request counter`, func() {
				go processIncrements(requestChan, invoiceChan)
				fetchInvoices(id, start, end, requestChan, invoiceChan)
				wg.Wait()

				So(requestCount, ShouldEqual, 1)
			})

			Reset(reset)
		})

	})
}
