#!/bin/sh

cd /usr/src
dlv debug --headless --listen=:2345 --api-version=2 --log ./cmd/main.go
