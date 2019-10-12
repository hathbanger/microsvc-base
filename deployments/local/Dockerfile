# golang
FROM golang

# build args
ARG SVC
ARG ENV
ARG CONSUL_TEMPLATE_RELEASE
ARG CONSUL_HTTP_ADDR
ARG CONSUL_HTTP_TOKEN
ARG VAULT_ADDR
ARG VAULT_TOKEN
ARG PORT

ENV SVC ${SVC}
ENV ENV ${ENV}
ENV CONSUL_HTTP_ADDR ${CONSUL_HTTP_ADDR}
ENV CONSUL_HTTP_TOKEN ${CONSUL_HTTP_TOKEN}
ENV CONSUL_HTTP_SSL_VERIFY=true
ENV DEBUG=false
ENV VAULT_ADDR ${VAULT_ADDR}
ENV VAULT_TOKEN ${VAULT_TOKEN}
ENV VAULT_SKIP_VERIFY=true
ENV PORT ${PORT}

# create the code location on the container
RUN mkdir -p /go/src/github.com/hathbanger/${SVC}

# copy the local code to the container location
COPY . /go/src/github.com/hathbanger/${SVC}

# install go-watcher
RUN go get github.com/canthefason/go-watcher

# install the binary under go/bin
RUN go install github.com/canthefason/go-watcher/cmd/watcher

# move the container's working directory
WORKDIR /go/src/github.com/hathbanger/${SVC}

# install consul-template
RUN curl -s https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_RELEASE}/consul-template_${CONSUL_TEMPLATE_RELEASE}_linux_amd64.tgz | tar -C /usr/local/bin -zx

ENTRYPOINT deployments/local/signal_handler.sh

