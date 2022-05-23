FROM golang:1.18 AS builder

COPY ../.. /src
WORKDIR /src

RUN cd app/system && GOPROXY=https://goproxy.cn CGO_ENABLED=0 make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

CMD ["./server"]
