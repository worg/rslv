FROM golang:alpine # Use the small size image from alpine
MAINTAINER worg <docker@worg.xyz>

WORKDIR /go/src/cmd/rslv
COPY cmd/rslv .

RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]