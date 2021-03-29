#build stage
FROM golang:1.15 AS builder
RUN apt-get install git
WORKDIR /go/src/app
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
COPY . /go/bin



#final stage
FROM ubuntu:latest
RUN apt-get update
RUN apt-get install ca-certificates -y
WORKDIR /app
RUN chmod -R 777 .
COPY --from=builder /go/bin /app
ENV GIN_MODE=release
ENTRYPOINT /app/mo2search
LABEL Name=mo2search Version=0.0.1
