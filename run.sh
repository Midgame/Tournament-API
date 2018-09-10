#!/bin/bash

go build
export PORT=5000
./Tournament-API -logtostderr=true
