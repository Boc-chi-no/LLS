FROM golang:1.20-alpine

WORKDIR /go/src/app
COPY . .
RUN go install github.com/rakyll/statik
RUN go generate
RUN go install linkshortener

CMD ["linkshortener"]