#!/bin/bash

go build
PORT=5000
./Tournament-API -logtostderr=true
