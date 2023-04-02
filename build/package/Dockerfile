FROM golang:1.20.0 AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /src/app

COPY . .

#dependencies
RUN go mod download
RUN go mod tidy

WORKDIR /src/app/cmd/app

RUN GOOS=linux go build -a -o app app.go

#Second stage for a smaller image
FROM alpine:3.17.2
RUN apk --no-cache add ca-certificates
WORKDIR /cmd/
COPY --from=builder /src/app/cmd/app/app .

ENV GIN_MODE "release"
CMD ["./app"]