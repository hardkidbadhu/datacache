FROM golang:1.13 AS builder

WORKDIR $GOPATH/src/

COPY ./datamodels ./datamodels
COPY ./workerservice ./workerservice

WORKDIR $GOPATH/src/workerservice

RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o worker_service main.go

FROM alpine

MAINTAINER badhmanabanm@hotmail.com

COPY --from=builder /go/src/workerservice/worker_service /worker_service

CMD ["./worker_service", "-conf", "./config.json"]