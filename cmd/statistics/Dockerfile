FROM golang:1

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -trimpath -race -o ./BUILD ./cmd/statistics

CMD ["./BUILD", "-p", "./cmd/statistics/config/local.toml", "-l", "debug"]