FROM golang:1.13 AS builder

WORKDIR $GOPATH/src/
COPY . ./datamodels
COPY . ./apiservice
RUN pwd
RUN ls

WORKDIR $GOPATH/src/apiservice
RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o apiservice /apiservice/main.go

FROM alpine
MAINTAINER badhmanabanm@hotmail.com

WORKDIR /datache

ADD config /datache/config

COPY --from=builder /datache .

CMD ["./apiservice"]
