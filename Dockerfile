FROM golang:1.18.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go mod download
RUN go build -o run-server

EXPOSE 8080

CMD [ "./run-server" ]