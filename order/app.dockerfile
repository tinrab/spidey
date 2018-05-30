FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/tinrab/spidey/order
COPY vendor ../vendor
COPY order ./
RUN go build -o /go/bin/app ./cmd/order/main.go

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
