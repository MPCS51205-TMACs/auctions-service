FROM golang:1.19
WORKDIR /go/src/auctions-service
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN mkdir -p /go/bin/auctions-service
RUN go build -v -o /go/bin/auctions-service ./...
WORKDIR /go/bin/auctions-service
CMD ["./main", "sql"]