#!/bin/bash

go get
go build
export PORT=5000
./Tournament-API -logtostderr=true
