FROM golang:1.25 AS builder

ARG SERVICE_DIR
WORKDIR /src

COPY go.work ./
COPY dogapm ./dogapm
COPY protos ./protos
COPY ordersvc ./ordersvc
COPY skusvc ./skusvc
COPY usrsvc ./usrsvc
COPY dogalarm ./dogalarm

RUN test -n "$SERVICE_DIR"
RUN cd "$SERVICE_DIR" && go build -o /out/app .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8080 8091 8092 8001 8002
ENTRYPOINT ["/app/app"]
