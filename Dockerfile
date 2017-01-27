FROM golang:1.7.4-wheezy

WORKDIR /go/src/github.com/SchweizerischeBundesbahnen/openshift-monitoring/deamon/

COPY ./deamon /go/src/github.com/SchweizerischeBundesbahnen/openshift-monitoring/deamon/

RUN go get github.com/cenkalti/rpc2
RUN go get github.com/gorilla/websocket
RUN go get github.com/mitchellh/mapstructure

RUN go install -v

CMD ["app"]