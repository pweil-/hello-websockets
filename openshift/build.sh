#!/bin/bash

GOPATH=$(cd .. && pwd)

pushd ../src
GOPATH=${GOPATH} go get golang.org/x/net/websocket
GOPATH=${GOPATH} go build server.go
mv server ../openshift/bin/.
popd

docker build -t pweil/hello-websocket .
