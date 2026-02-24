FROM golang:1.25-alpine
WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o my-exporter main.go

EXPOSE 2112
CMD ["/app/my-exporter"]