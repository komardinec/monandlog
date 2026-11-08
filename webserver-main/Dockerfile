FROM golang@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS builder 
# 1.24-alpine3.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o webserver ./cmd/web

FROM alpine@sha256:93b41ef40b1677fce8fcc64ff4af09daddfacf5040afde618bc46c1e2d1cc7d4
# 3.24.4

WORKDIR /webserver/

RUN adduser -D webserver
USER webserver

COPY --from=builder /app/webserver .
COPY --from=builder /app/ui .

ENTRYPOINT ["./webserver"]
CMD ["-addr=:8080"]
