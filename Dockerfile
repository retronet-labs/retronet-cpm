FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY . .

RUN go build -o /out/retronet-cpm ./cmd/retronet-cpm \
    && go build -o /out/retronet-cpm-live ./cmd/retronet-cpm-live

FROM alpine:latest

WORKDIR /app
COPY --from=builder /out/retronet-cpm /app/retronet-cpm
COPY --from=builder /out/retronet-cpm-live /app/retronet-cpm-live

ENTRYPOINT ["/app/retronet-cpm"]
CMD ["-h"]
