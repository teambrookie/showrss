#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o showrss .
docker build -t showrss .