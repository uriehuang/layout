FROM golang:1.24 AS builder

COPY . /src
WORKDIR /src

ARG GITHUB_USER
ARG GITHUB_TOKEN

ENV GITHUB_USER=$GITHUB_USER
ENV GITHUB_TOKEN=$GITHUB_TOKEN

RUN export GOPRIVATE="github.com/uriehuang/*" \
    && git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf https://github.com \
    && make build

FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app/bin
COPY --from=builder /src/configs /app/configs
COPY start.sh .

ARG CONFIG_FILE_NAME
ENV CONFIG_FILE_NAME=$CONFIG_FILE_NAME

RUN chmod +x start.sh

EXPOSE 8000
EXPOSE 9000

CMD ["sh", "-c", "./start.sh ${CONFIG_FILE_NAME}"]
