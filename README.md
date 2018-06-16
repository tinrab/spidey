# Spidey

Online store based on microservices and GraphQL.

The underlying source code for the article [Using GraphQL API Gateway for Microservices in Go](https://outcrawl.com/go-graphql-gateway-microservices).

## Build

```
$ vgo mod -vendor
$ docker-compose up -d --build
```

Open <http://localhost:8000/playground> in your browser.
