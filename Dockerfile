FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server

FROM scratch

COPY --from=builder /server /server

EXPOSE 50051

CMD ["/server"]

