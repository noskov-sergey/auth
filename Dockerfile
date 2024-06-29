FROM golang:1.20.3-alpine AS builder

COPY . /github.com/noskov-sergey/auth/source/
WORKDIR /github.com/noskov-sergey/auth/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/noskov-sergey/auth/source/bin/crud_server .

CMD ["./crud_server"]