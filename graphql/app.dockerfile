FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/tinrab/spidey/graphql
COPY vendor ../vendor
COPY account ../account
COPY catalog ../catalog
COPY order ../order
COPY graphql ./
RUN go build -o /go/bin/app

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
