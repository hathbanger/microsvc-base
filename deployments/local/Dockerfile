# golang
FROM golang

ARG SVC=microsvc-base
ARG PORT
ARG ENV=dev

ENV SVC ${SVC}
ENV ENV ${ENV}
ENV PORT ${PORT}

RUN mkdir -p /go/src/github.com/hathbanger/${SVC}
COPY . /go/src/github.com/hathbanger/${SVC}
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher
WORKDIR /go/src/github.com/hathbanger/${SVC}
RUN echo "${PORT}"

ENTRYPOINT watcher -run github.com/hathbanger/${SVC}/cmd -config /go/src/github.com/hathbanger/${SVC}/configs/config.json -port ${PORT} -watch github.com/hathbanger/${SVC}
