FROM golang:1.21.3-alpine3.18

WORKDIR /usr/src/api

COPY ./ .

RUN go build ./cmd/api.go

CMD ["./api"]
