package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync/atomic"
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
	// ErrorApiLimit is returned when we must stop making calls
	ErrorAPILimit = `API limit reached`
	// API base url
	baseURL       = `http://34.209.24.195/facturas`
	requestFormat = baseURL + `?id=%s&start=%s&finish=%s`
	// Time format allowed
	timeFmt = `2006-01-02`
)

var (
	startDate, endDate, id string
	debug                  bool
	// ErrorExceededCount is returned when API returns a `more than… found` message
	ErrorExceededCount = errors.New(`API limit reached`)
)

// Job holds pending request data
type job struct {
	id         string
	start, end time.Time
}

func init() {
	// I use golang's flag to parse command line options
	flag.StringVar(&startDate, `start`, ``, `Start of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&endDate, `finish`, ``, `End of date range to find invoices [YYYY-MM-DD]`)
	flag.StringVar(&id, `id`, ``, `User id to fetch invoices`)
	flag.Parse()

	// fallback load data from ENV vars
	if startDate == `` {
		os.Getenv(`START_DATE`)
	}

	if endDate == `` {
		os.Getenv(`END_DATE`)
	}

	if id == `` {
		os.Getenv(`USER_ID`)
	}

	debug = os.Getenv(`RSLV_DEBUG`) != ``
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

	c := make(chan job, 1)
	c <- job{
		id:    id,
		start: start,
		end:   end,
	}

	invoiceCount, requestCount := processData(c)

	fmt.Printf("%d invoices were found, using %d requests\n", invoiceCount, requestCount)
}

func processData(c chan job) (int, int) {
	done := make(chan struct{})
	var invoiceCount int
	// we'll use int32 to keep an atomic count
	// this is intended to prevent race conditions
	var doneCount, requestCount int32
	var stop = false // stop flag [as we don't know the exact number of iterations]

	defer func() {
		// close channels at the end
		close(done)
		close(c)
	}()

	for !stop {
		select {
		case j, ok := <-c:
			if !ok {
				stop = true
			}

			rc := atomic.LoadInt32(&requestCount)
			rc++
			atomic.StoreInt32(&requestCount, rc)

			count, err := FetchInvoices(j.id, j.start, j.end)
			go func() {
				time.Sleep(1)
				done <- struct{}{}
			}()

			if err != nil {
				// split job into two new ranges
				go SplitJob(c, j)
				continue
			}
			invoiceCount += count

		case <-done:
			dc := atomic.LoadInt32(&doneCount)
			rc := atomic.LoadInt32(&requestCount)
			atomic.StoreInt32(&doneCount, dc+1)
			dc++
			if rc == dc {
				stop = true
			}
		}
	}

	return invoiceCount, int(requestCount)
}

// SplitJob takes a job and creates two new from
// half ranges of the original job date or a single
// new job when dates are too close
func SplitJob(c chan job, j job) {
	start := j.start
	end := j.end
	id := j.id

	daysBetween := GetDaysBetween(start, end)

	// do we have a very short time span?
	if daysBetween <= 3 {
		c <- job{
			id:    id,
			start: AddDays(start, 1),
			end:   end,
		}
		return
	}

	half := daysBetween / 2
	c <- job{
		id:    id,
		start: start,
		end:   AddDays(start, half),
	}

	c <- job{
		id:    id,
		start: AddDays(end, -half),
		end:   end,
	}
}

// GetDaysBetween returns the days elapsed within two dates
func GetDaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

// AddDays returns a date with n days added [may be negative]
func AddDays(date time.Time, days int) time.Time {
	return date.Add(time.Hour * time.Duration(24*days))
}
