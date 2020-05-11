FROM golang:1.13 AS builder

WORKDIR $GOPATH/src/

COPY ./datamodels ./datamodels
COPY ./apiservice ./apiservice

WORKDIR $GOPATH/src/apiservice

RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o apiservice main.go

FROM alpine

MAINTAINER badhmanabanm@hotmail.com

COPY --from=builder /go/src/apiservice/apiservice /apiservice

CMD ["./apiservice", "-conf", "./config.json"]