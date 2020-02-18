FROM golang:latest

WORKDIR $GOPATH/src/go-blog
COPY . $GOPATH/src/go-blog
RUN go build -o go-blog main.go

EXPOSE 8080
ENTRYPOINT ["./go-blog"]