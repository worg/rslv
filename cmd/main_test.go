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

		/*
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
		*/

		Convey(`Utility functions`, func() {
			date, _ := time.Parse(timeFmt, `2017-01-01`)

			Convey(`AddDays should increment a date day count when receiving a positive number`, func() {
				resultDate, _ := time.Parse(timeFmt, `2017-01-04`)
				So(AddDays(date, 3), ShouldEqual, resultDate)
			})

			Convey(`AddDays should decrement a date when receiving a negative number`, func() {
				resultDate, _ := time.Parse(timeFmt, `2016-12-29`)
				So(AddDays(date, -3), ShouldEqual, resultDate)
			})

			Convey(`GetDaysBetween should return the number of days elapsed in two specified dates`, func() {
				endDate, _ := time.Parse(timeFmt, `2017-01-30`)
				So(GetDaysBetween(date, endDate), ShouldEqual, 29)
			})

			Convey(`SplitJob should send two jobs with the daterange in halves`, func() {
				c := make(chan job)
				endDate, _ := time.Parse(timeFmt, `2017-01-30`)

				go SplitJob(c, job{id: `1`, start: date, end: endDate})
				first := <-c
				second := <-c

				So(first.start, ShouldEqual, date)
				So(first.end, ShouldEqual, AddDays(endDate, -15))

				So(second.start, ShouldEqual, AddDays(date, 15))
				So(second.end, ShouldEqual, endDate)
			})

			Convey(`SplitJob should send a single job when dates are too close`, func() {
				c := make(chan job)
				endDate, _ := time.Parse(timeFmt, `2017-01-03`)

				go SplitJob(c, job{id: `1`, start: date, end: endDate})
				first := <-c

				So(first.start, ShouldEqual, AddDays(date, 1))
				So(first.end, ShouldEqual, endDate)
			})
		})

	})
}
