FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o 1337b04rd ./main.go

CMD ["./1337b04rd"]
