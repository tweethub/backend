FROM golang:1

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -trimpath -race -o ./BUILD ./cmd/tweets

CMD ["./BUILD", "-p", "./cmd/tweets/config/local.toml", "-l", "debug"]