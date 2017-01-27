FROM golang:1.7.4-wheezy

WORKDIR /go/src/

COPY ./deamon /go/src/
RUN ls .

RUN go get github.com/cenkalti/rpc2
RUN go get github.com/gorilla/websocket
RUN go get github.com/mitchellh/mapstructure

RUN go install -v

CMD ["app"]