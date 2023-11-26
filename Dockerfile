FROM golang:1.20-alpine AS build-env

WORKDIR /usr/src/app

# trying to cache modules independent of source code changes
COPY go.mod go.sum  ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p bin
RUN go build -tags musl -o bin ./...

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build-env /usr/src/app/bin .

EXPOSE 80 81
CMD ["./ketalk-api"]