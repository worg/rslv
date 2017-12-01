package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
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
			id = `1`

			endDate = `2017 Jan 03`
			So(main, ShouldPanicWith, ErrorInvalidFormat)

			startDate = `2017 Jan 01`
			So(main, ShouldPanicWith, ErrorInvalidFormat)
		})

		Convey(`When making requests`, func() {
			startDate = `2017-01-01`
			endDate = `2017-03-01`
			id = `1`
			start, _ := time.Parse(timeFmt, startDate)
			end, _ := time.Parse(timeFmt, endDate)

			httpmock.Activate() // Enable http mocker

			Convey(`FetchInvoices should return the invoice count on valid response`, func() {
				// Fake the http response
				httpmock.RegisterResponder(`GET`, baseURL, httpmock.NewStringResponder(200, `40`))

				c, err := FetchInvoices(id, start, end)
				So(err, ShouldEqual, nil)
				So(c, ShouldEqual, 40)

				Reset(httpmock.Reset)
			})

			Convey(`FetchInvoices should handle a truncated results response`, func() {
				apiCalls := 0
				// Fake the http response
				httpmock.RegisterResponder(`GET`, baseURL, func(r *http.Request) (*http.Response, error) {
					apiCalls++

					// Fail when apiCalls is < 2
					if apiCalls < 2 {
						return httpmock.NewStringResponse(200, `Hay más de 100 resultados`), nil
					}

					return httpmock.NewStringResponse(200, fmt.Sprint(apiCalls+1)), nil
				})

				c, err := FetchInvoices(id, start, end)
				So(err, ShouldEqual, ErrorExceededCount)
				So(c, ShouldEqual, 0)

				c, err = FetchInvoices(id, start, end)
				So(err, ShouldEqual, nil)
				So(c, ShouldEqual, 3)

				Reset(httpmock.Reset)
			})

			Convey(`FetchInvoices should handle API rate limits`, func() {
				httpmock.RegisterResponder(`GET`, baseURL, httpmock.NewStringResponder(400, `"API limit reached"`))

				testFn := func() {
					FetchInvoices(id, start, end)
				}

				So(testFn, ShouldPanicWith, ErrorAPILimit)

				Reset(httpmock.Reset)
			})

			Convey(`FetchInvoices should handle unexpected errors`, func() {
				httpmock.RegisterResponder(`GET`, baseURL, httpmock.NewStringResponder(500, `"LOL"`))

				testFn := func() {
					FetchInvoices(id, start, end)
				}

				So(testFn, ShouldNotPanicWith, ErrorAPILimit)

				Reset(httpmock.Reset)
			})

			Reset(httpmock.DeactivateAndReset)
		})

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
