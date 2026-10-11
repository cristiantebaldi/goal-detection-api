#!/bin/sh
cd "$1"
cp config.example.yaml config.yaml
sed -i -e 's/$SERVER_PORT/'8080'/g' config.yaml
