FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY resume.pdf resume.png ./
COPY 2-page/resume.pdf ./2-page/resume.pdf

EXPOSE $PORT

CMD ["./server"]
