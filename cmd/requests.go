package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// FetchInvoices gets the invoice count [or error] for a particular time span
// returns error when API fails
func FetchInvoices(id string, start, end time.Time) (int, error) {
	startString := start.Format(timeFmt)
	endString := end.Format(timeFmt)

	// Do the API call
	resp, err := http.Get(fmt.Sprintf(requestFormat, id, startString, endString))
	if err != nil {
		panic(err) // we treat this errors as fatal
	}

	// Get the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // we treat this errors as fatal
	}

	// try to get the count from the response body
	count, err := strconv.Atoi(string(body))
	var e *strconv.NumError
	if err != nil {
		e = err.(*strconv.NumError)
		if e.Err != strconv.ErrSyntax {
			panic(err)
		}

		return 0, ErrorExceededCount
	}

	return count, nil
}
