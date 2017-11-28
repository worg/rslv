package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResuelve(t *testing.T) {
	Convey(`On running the app`, t, func() {
		Convey(`Input must not be empty`, func() {
			So(main, ShouldPanicWith, ErrorInvalidInput)
		})

		Convey(`Dates must differ`, func() {
			startDate = `2017-11-01`
			endDate = `2017-11-01`
			id = `1`
			So(main, ShouldPanicWith, ErrorInvalidRange)
		})

	})
}
