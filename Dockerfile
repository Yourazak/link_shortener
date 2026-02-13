FROM golang:1.25.0-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o link_shortener ./cmd/link_shortener/

CMD [ "./link_shortener" ]