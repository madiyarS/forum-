#!/bin/sh
docker build -t forum-app .
docker run -d -p 8082:8080 --name forum -v "$(pwd -W)/data:/root/data" forum-app