#!/bin/bash


# build
cd src && \
env \
    GOOS=linux \
    GOARCH=arm \
    GOARM=5 \
go build -o ../build/tempsens.bin .