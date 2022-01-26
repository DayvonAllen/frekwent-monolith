FROM golang:latest

LABEL maintainer="dev <test@test.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build

CMD ["./freq"]