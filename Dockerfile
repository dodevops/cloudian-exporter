FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o cloudian-exporter cmd/cloudian-exporter.go

FROM alpine

COPY --from=builder /app/cloudian-exporter /
RUN adduser -D cloudian-exporter && chmod +x /cloudian-exporter

USER 1001
EXPOSE 8080
ENTRYPOINT ["/cloudian-exporter"]