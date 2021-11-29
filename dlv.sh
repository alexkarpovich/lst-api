#!/bin/sh

cd /go/src/yuzik/api
dlv debug --headless --listen=:2345 --api-version=2 --log main.go
