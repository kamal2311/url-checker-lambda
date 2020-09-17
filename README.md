# url-checker-lambda

This lambda function can be set as an Http proxy trigger behind an 
AWS API Gateway.

It takes an Http request with a URL and responds back with information about the URL
for its inclusion in malware feeds.

Http GET `/v1/url-info/{URL to test}` responds with
 ```json
{
    "is_safe": false,
    "score": 8,
    "source": "Sophos"
}
```

OR 

```json
{
  "is_safe": true
}
```

## Run tests ( Needs a dynamodb table in your AWS account)
export AWS_PROFILE={{your AWS profile}}
`go test -v www.github.com/kamal2311/url-checker-lambda`

## Deployment scripts
Refer to `/scripts` folder
