export AWS_PROFILE=kamal
aws lambda create-function --function-name my-function --runtime go1.x \
  --zip-file fileb://main.zip --handler main \
  --role arn:aws:iam::132491518201:role/service-role/basic-lambda-role