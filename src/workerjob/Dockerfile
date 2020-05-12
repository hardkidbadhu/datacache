FROM golang:1.13 AS builder

WORKDIR $GOPATH/src/

COPY ./datamodels ./datamodels
COPY ./workerjob ./workerjob

WORKDIR $GOPATH/src/workerjob

RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o workerjob main.go

FROM alpine

MAINTAINER badhmanabanm@hotmail.com

COPY --from=builder /go/src/workerjob/workerjob /workerjob

CMD ["./workerjob", "-conf", "./config.json"]