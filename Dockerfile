FROM golang:alpine

WORKDIR /go/src/app
COPY .. .
RUN go get github.com/kataras/iris/v12@master && \
	go build

CMD ["./web-text-manager"]
