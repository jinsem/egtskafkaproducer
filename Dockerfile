FROM golang:1.12.7

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o egtsreceiver ./cmd/egtsreceiver

CMD ["./egtsreceiver"]