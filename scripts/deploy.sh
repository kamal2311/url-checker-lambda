#!/usr/bin/env bash
set -eo pipefail

LAMBDA_NAME=url-checker-lambda

# Create the lambda execution role
aws iam create-role --role-name ${LAMBDA_NAME}-role --assume-role-policy-document file://trust-policy.json
echo "Created the role"

# Attach policies
aws iam attach-role-policy --role-name ${LAMBDA_NAME}-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
echo "Attached the policy to the role"

# Create lambda function
aws lambda create-function --function-name ${LAMBDA_NAME} --runtime go1.x \
  --zip-file fileb://main.zip --handler main \
  --role arn:aws:iam::132491518201:role/${LAMBDA_NAME}-role

echo "Created the lambda function"

MC_TABLE_NAME=malicious-urls

# Create a dynamo-db table
aws dynamodb create-table \
    --table-name $MC_TABLE_NAME \
    --attribute-definitions \
        AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
fi

# Insert some items
aws dynamodb put-item \
--table-name $MC_TABLE_NAME  \
--item \
    '{"Id": {"S": "66b9cb08638d49a6d3559718551d59243fa2b0eb"}, "Score": {"N": "8"}, "Source": {"S": "Sophos"}}' \
--return-consumed-capacity TOTAL

# Insert some items
aws dynamodb put-item \
--table-name $MC_TABLE_NAME  \
--item \
    '{"Id": {"S": "4782cc39a5294f566242f9d36bccc9889916e2b6"}, "Score": {"N": "9"}, "Source": {"S": "Malware Patrol"}}' \
--return-consumed-capacity TOTAL
