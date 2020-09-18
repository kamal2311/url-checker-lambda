#!/usr/bin/env bash
set -eo pipefail

# Test table name
export MC_TABLE_NAME=test-malicious-urls

# Create a new table if it does not exist already
if [ -z $(aws dynamodb describe-table --table-name ${MC_TABLE_NAME})]; then
aws dynamodb create-table \
    --table-name $MC_TABLE_NAME \
    --attribute-definitions \
        AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
fi

# Add an item
aws dynamodb put-item \
--table-name $MC_TABLE_NAME  \
--item \
    '{"Id": {"S": "66b9cb08638d49a6d3559718551d59243fa2b0eb"}, "Score": {"N": "8"}, "Source": {"S": "Sophos"}}' \
--return-consumed-capacity TOTAL

# Add an item
aws dynamodb put-item \
--table-name $MC_TABLE_NAME  \
--item \
    '{"Id": {"S": "4782cc39a5294f566242f9d36bccc9889916e2b6"}, "Score": {"N": "9"}, "Source": {"S": "Malware Patrol"}}' \
--return-consumed-capacity TOTAL

# Run service and data tests
go test -v www.github.com/kamal2311/url-checker-lambda

# Delete the test table
aws dynamodb delete-table --table-name $MC_TABLE_NAME