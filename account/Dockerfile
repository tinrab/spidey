FROM golang:1.10.2 AS build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN env GOOS=linux GARCH=amd64 go install -v -a -tags netgo -installsuffix netgo -ldflags "-linkmode external -extldflags -static"

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD ["app"]
