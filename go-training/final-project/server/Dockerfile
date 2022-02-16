FROM golang:1.16-alpine
WORKDIR /app
COPY main.go /app
COPY go.mod /app
COPY data/data.json /app
RUN go build -o server
EXPOSE 8080
CMD ["./server"]