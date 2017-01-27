FROM golang:1.7.4-onbuild

COPY ./deamon /go/src/app

#RUN go get github.com/cenkalti/rpc2 &&
#    go get github.com/gorilla/websocket &&
#    go get github.com/mitchellh/mapstructure
