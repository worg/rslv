package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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

			reset := func() {
				requestCount = 0
				invoiceCount = 0
			}
			reset()

			Convey(`Counters should start on zero`, func() {
				So(requestCount, ShouldEqual, 0)
				So(invoiceCount, ShouldEqual, 0)
			})

			Reset(reset)
		})

	})
}
