FROM golang:1.8-jessie

WORKDIR /go/src/github.com/oscp/openshift-monitoring/daemon/

COPY . /go/src/github.com/oscp/openshift-monitoring/

RUN go get github.com/cenkalti/rpc2
RUN go get github.com/gorilla/websocket
RUN go get github.com/mitchellh/mapstructure

RUN go install -v

# Install necessary tools
RUN apt-get update && apt-get install -y --no-install-recommends dnsutils

CMD ["daemon"]