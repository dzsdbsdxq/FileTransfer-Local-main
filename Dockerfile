FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/FileTransfer
COPY . $GOPATH/src/FileTransfer
RUN go build .
EXPOSE 8999
ENTRYPOINT ["./FileTransfer"]