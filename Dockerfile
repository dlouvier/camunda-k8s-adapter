FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/dlouvier/camunda-k8s-adapter
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/github.com/dlouvier/camunda-k8s-adapter/app .
CMD ["./app"]