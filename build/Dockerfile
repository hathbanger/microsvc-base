# build the binaries
FROM golang:1.11.5 AS builder
ARG RELEASE_DIR
ARG SERVICE_NAME
RUN mkdir -p {$RELEASE_DIR,/go/src/github.com/hathbanger/$SERVICE_NAME}
CMD GO111MODULE=on
COPY . /go/src/github.com/hathbanger/$SERVICE_NAME
WORKDIR /go/src/github.com/hathbanger/$SERVICE_NAME
RUN make build

FROM busybox AS service
ARG RELEASE_DIR
ARG SERVICE_NAME
ARG SERVICE_DIR
COPY --from=builder $RELEASE_DIR/ $SERVICE_DIR/
COPY --from=builder /go/src/github.com/hathbanger/$SERVICE_NAME/configs/config.json $SERVICE_DIR/config.json
RUN chmod +x $SERVICE_DIR/$SERVICE_NAME && \
  ls -lh $SERVICE_DIR

FROM service
ARG SERVICE_NAME
ARG SERVICE_DIR
ARG SERVICE_PORT
ENV PORT=$SERVICE_PORT
ENV SERVICE=$SERVICE_NAME
WORKDIR $SERVICE_DIR
RUN ls
ENTRYPOINT ./$SERVICE
