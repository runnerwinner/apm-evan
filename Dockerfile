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
ENV TZ=Asia/Shanghai
RUN apt-get update \
	&& DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates curl gnupg tzdata \
	&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
	&& echo $TZ > /etc/timezone \
	&& curl -fsSL https://artifacts.elastic.co/GPG-KEY-elasticsearch | gpg --dearmor -o /usr/share/keyrings/elastic-keyring.gpg \
	&& echo "deb [signed-by=/usr/share/keyrings/elastic-keyring.gpg] https://artifacts.elastic.co/packages/7.x/apt stable main" > /etc/apt/sources.list.d/elastic-7.x.list \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends filebeat \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/app /app/app
COPY conf/filebeat.yml /etc/filebeat/filebeat.yml
COPY conf/start-with-filebeat.sh /usr/local/bin/start-with-filebeat.sh
RUN chmod +x /usr/local/bin/start-with-filebeat.sh
EXPOSE 8080 8091 8092 8001 8002
ENTRYPOINT ["/usr/local/bin/start-with-filebeat.sh"]
