FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY . .

RUN go build -o /out/retronet-cpm ./cmd/retronet-cpm

FROM alpine:latest

WORKDIR /app
COPY --from=builder /out/retronet-cpm /app/retronet-cpm

ENTRYPOINT ["/app/retronet-cpm"]
CMD ["-h"]
