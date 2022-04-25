##
## Build
##
FROM golang:1.18.1-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /homework-microservice .

##
## Deploy
##
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=build /homework-microservice ./homework-microservice

RUN chmod +x ./homework-microservice

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["./homework-microservice"]