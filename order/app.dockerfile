FROM golang:1.11.0-alpine3.8 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/tinrab/spidey/order
COPY vendor ../vendor
COPY account ../account
COPY catalog ../catalog
COPY order ./
RUN go build -o /go/bin/app ./cmd/order/main.go

FROM alpine:3.8
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
