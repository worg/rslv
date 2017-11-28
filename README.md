## Problem

We need to fetch the total invoice count within a date span from a limited API, the API returns `more than 100 results` if the query exceedes that count.


## Solution

I used a _divide and conquer_ algorithm, if the first call exceedes the result count, I'll split the date in two ranges, and make two calls, one for each range, This process repeats until we get the total count.


## Running

TBD
