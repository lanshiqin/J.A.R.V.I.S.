FROM golang:1.12

ADD ./ /go/src/jarvis/

RUN cd /go/src/jarvis/system/core && go run main.go