# build the binaries and configs
FROM golang:1.12.1 AS builder
ARG RELEASE_DIR
ARG SERVICE_NAME
ARG VAULT_RELEASE
ARG CONSUL_TEMPLATE_RELEASE=0.19.5
ARG PORT
RUN mkdir -p {$RELEASE_DIR,/go/src/github.com/hathbanger/$SERVICE_NAME}
COPY . /go/src/github.com/hathbanger/$SERVICE_NAME
WORKDIR /go/src/github.com/hathbanger/$SERVICE_NAME
RUN apt-get update && apt-get install unzip && make build && \
   curl -s https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_RELEASE}/consul-template_${CONSUL_TEMPLATE_RELEASE}_linux_amd64.tgz | tar -C /usr/local/bin -zvx && \
   curl -sO https://releases.hashicorp.com/vault/${VAULT_RELEASE}/vault_${VAULT_RELEASE}_linux_amd64.zip; unzip -d /usr/local/bin vault_${VAULT_RELEASE}_linux_amd64.zip; rm vault_${VAULT_RELEASE}_linux_amd64.zip; \
   ls -lh $RELEASE_DIR

# grab the binaries from builder
FROM alpine AS config
ARG RELEASE_DIR
ARG SERVICE_NAME
ARG SERVICE_DIR
COPY --from=builder /go/src/github.com/atu/${SERVICE_NAME}/$RELEASE_DIR/ $SERVICE_DIR/
COPY --from=builder /go/src/github.com/atu/${SERVICE_NAME}/configs $SERVICE_DIR/configs/
COPY --from=builder /go/src/github.com/atu/${SERVICE_NAME}/tools/signal_handler.sh $SERVICE_DIR/
COPY --from=builder /go/src/github.com/atu/${SERVICE_NAME}/Makefile $SERVICE_DIR/
RUN mkdir -p $SERVICE_DIR && mv $SERVICE_DIR/$SERVICE_NAME-* $SERVICE_DIR/$SERVICE_NAME && \
    chmod +x $SERVICE_DIR/$SERVICE_NAME && \
    ls -lh $SERVICE_DIR

# stage the service contianer
FROM alpine AS service
ARG SERVICE_NAME
ARG SERVICE_DIR
ARG GITHUB_TOKEN
COPY --from=builder /usr/local/bin/consul-template /usr/local/bin/
COPY --from=builder /usr/local/bin/vault /usr/local/bin/
COPY --from=builder /go/src/github.com/hathbanger/${SERVICE_NAME}/tools $SERVICE_DIR/tools/
COPY --from=builder /go/src/github.com/hathbanger/${SERVICE_NAME}/tools/etc/rsyslog.d/00-elasticsearch.conf /etc/rsyslog.d/
COPY --from=config $SERVICE_DIR $SERVICE_DIR
RUN apk add --no-cache \
    make ca-certificates rsyslog rsyslog-elasticsearch rsyslog-mmjsonparse \
    && ls -ltrh $SERVICE_DIR

# final service container
FROM service
ARG SERVICE_DIR
ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=$GITHUB_TOKEN
ENV PORT $PORT
EXPOSE $SERVICE_PORT
WORKDIR $SERVICE_DIR
ENTRYPOINT ./signal_handler.sh
