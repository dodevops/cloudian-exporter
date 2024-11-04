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

# Labels

LABEL io.artifacthub.package.readme-url=https://github.com/dodevops/cloudian-exporter
LABEL org.opencontainers.image.description="This Prometheus exporter provides a metrics endpoint with metrics from a Cloudian installation."
LABEL org.opencontainers.image.documentation=https://github.com/dodevops/cloudian-exporter
LABEL org.opencontainers.image.source=https://github.com/dodevops/cloudian-exporter
LABEL org.opencontainers.image.title="cloudian-exporter"
LABEL org.opencontainers.image.url=https://github.com/dodevops/cloudian-exporter
LABEL org.opencontainers.image.vendor="DO! DevOps"