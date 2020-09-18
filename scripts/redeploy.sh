#!/usr/bin/env bash
set -eo pipefail

# Build the go binary
GOOS=linux go build -o main www.github.com/kamal2311/url-checker-lambda

## Package the binary into a zip file
# For windows
~/go/bin/build-lambda-zip.exe -output main.zip main
# For Linux/Mac, uncomment the following line and comment the above line
# zip main.zip main

## Update lambda function with from the zip above
aws lambda update-function-code --function-name my-function --zip-file fileb://main.zip
