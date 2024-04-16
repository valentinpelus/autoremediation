##
## Build phase
##
FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/docker.io/library/ubuntu:22.04 AS build

LABEL maintainer="Bedrock"

RUN apt-get update && \
    apt-get install curl -y
RUN apt-get install -y ca-certificates curl
RUN apt-get install -y apt-transport-https
    #curl -L https://github.m6web.fr/valentin-pelus/autoremediate/releases/download/v.0.5.0-alpha/remediate-bsd-hapee4-linux-amd64 -o /app/remediate-bsd-hapee && \
    #chmod +x /app/remediate-bsd-hapee
COPY auto-remediation-linux-amd64 /remediate-bsd-hapee
RUN mkdir /var/src
COPY . /var/src/
RUN    chmod +x /remediate-bsd-hapee

##
## Deploy
##

#FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/gcr.io/distroless/base
FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/docker.io/library/ubuntu:22.04

WORKDIR /app

COPY --from=build /remediate-bsd-hapee /app/autoremediate
COPY --from=build /var/src /tmp/src/
COPY config.yaml /app/config.yaml
ENTRYPOINT /app/autoremediate

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y apt-transport-https

ENV CONFIG_PATH "/config.yaml"

CMD ["/app/autoremediate"]

ARG BUILD_DATE
ARG VCS_TYPE=git
ARG VCS_URL
ARG VCS_REF
LABEL build-date=$BUILD_DATE \
    vcs-type=$VCS_TYPE \
    vcs-url=$VCS_URL \
    vcs-ref=$VCS_REF
