##
## Build
##
FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/docker.io/library/ubuntu:20.04 AS build
LABEL maintainer="Bedrock"

RUN apt-get update && \
    apt-get install curl -y && \
    mkdir /app && \
    curl -L https://github.com/BedrockStreaming/alerts-validator/releases/download/v2.1.3/alerts-validator-linux-amd64 -o /app/alerts-validator && \
    chmod +x /app/alerts-validator

##
## Deploy
##
FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/gcr.io/distroless/base

WORKDIR /app

COPY --from=build /app/alerts-validator /app/alerts-validator

ENV CONFIG_PATH ""

USER nonroot:nonroot

CMD ["/app/alerts-validator"]

ARG BUILD_DATE
ARG VCS_TYPE=git
ARG VCS_URL
ARG VCS_REF
LABEL build-date=$BUILD_DATE \
    vcs-type=$VCS_TYPE \
    vcs-url=$VCS_URL \
    vcs-ref=$VCS_REF
