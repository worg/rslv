# rslv - Resuelve invoice tool

[![Build Status](https://img.shields.io/travis/worg/rslv/master.svg?style=flat-square)](https://travis-ci.org/worg/rslv)
[![Go Report Card](https://goreportcard.com/badge/github.com/worg/rslv/cmd?style=flat-square)](https://goreportcard.com/report/github.com/worg/rslv)


## Problem

We need to fetch the total invoice count within a date span from a limited API, the API returns `more than 100 results` if the query exceedes that count.


## Solution

I used a _divide and conquer_ algorithm, if the first call exceedes the result count, I'll split the date in two ranges, and make two calls, one for each range, This process repeats until we get the total count.


## Running

I choosed docker to deploy this utility. for this makes easier to not depend on the user environment for running this program.
I'll assume the user has docker in her/his system.

### Options

`rslv` takes three options, either as command line flags or environment variables:

| Option | Environment Variable | Command Line Flag | Description                          | Format       |
| --     | :--:                 | :--:              | :--:                                 | :--:         |
| ID     | USER_ID              | id                | User id to fetch invoices            | GUID         |
| Start  | START_DATE           | finish            | End of date range to find invoices   | [YYYY-MM-DD] |
| Finish | END_DATE             | start             | Start of date range to find invoices | [YYYY-MM-DD] |


### Building the image

In order to have an usable base image with `rslv` and its dependencies we'll build an image containing golang and our dependencies.

* Clone the repository.
* Navigate to the repository root.
* Execute `docker build -t {TAG_NAME} .` and wait for our image to be built.
  *  Where `{TAG_NAME}` is a name for our image, for example: `rslv_img`

This will fetch our tool dependencies and build a binary

### Running `rslv`

Once built the image we can run our tool

Execute 

`docker run  --rm --name rslv {TAG_NAME} -id={USER_ID} -start={START_DATE} -finish={END_DATE}`

Replacing:

* `{TAG_NAME}` with the tag defined in the build step
* `{USER_ID}` with the user id to find invoices
* `{START_DATE}` with the start of date range within invoices must be searched
* `{END_DATE}` with the end of the date range to perfotm the search

I prefer to use `--rm` option to discard our container after execution is done.

### Output

At the end of execution `rslv` will output the total invoice count and how many requests it took.

`{X} invoices were found, using {Y} requests`

### Debugging 

Executing with `RSLV_DEBUG` environment variable distinct of empty string will output the requests made and their results, in the following format: 

```
-- REQUEST: "http://U.R.L/path?query=args"
 - RESPONSE: "HTTP_RESPONSE"
```

## Testing

Test results are available on travis-ci [see badge], alternatively one can run locally the tests.
Assuming a working golang setup:

* `cd cmd/rslv` from the repository root
* `go get -t ./...`  to fetch the project dependencies
* `go test -v -cover` to display the test results using the command line
* `goconvey -cover . ` to display test results using the browser



