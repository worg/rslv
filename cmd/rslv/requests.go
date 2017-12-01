package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FetchInvoices gets the invoice count [or error] for a particular time span
// returns error when API fails
func FetchInvoices(id string, start, end time.Time) (int, error) {
	startString := start.Format(timeFmt)
	endString := end.Format(timeFmt)

	// Do the API call
	url := fmt.Sprintf(requestFormat, id, startString, endString)
	resp, err := http.Get(url)
	if err != nil {
		panic(err) // we treat this error as fatal
	}

	// Get the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // we treat this error as fatal
	}

	stringBody := strings.Replace(string(body), `"`, ``, -1) // remove quotes from response body

	// fail on HTTP errors
	if resp.StatusCode > 200 {
		if stringBody == ErrorAPILimit {
			panic(ErrorAPILimit)
		}

		panic(stringBody)
	}

	if debug {
		// we spawn a goroutine to prevent IO blocking
		go fmt.Printf("-- REQUEST: %q \n - RESPONSE: %q \n\n", url, stringBody)
	}
	// try to get the count from the response body
	count, err := strconv.Atoi(stringBody)
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
