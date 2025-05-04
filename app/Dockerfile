FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

EXPOSE 1323

CMD ["go", "run", "main.go"]
