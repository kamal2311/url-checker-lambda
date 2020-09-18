# url-checker-lambda

This lambda function can be set as an Http proxy trigger behind an 
AWS API Gateway.

It takes an Http request with a URL and responds back with information about the URL
for its inclusion in malware feeds.

## Test drive a deployed instance of the application

**NOTE** :
This API is not publicly available and will require you to supply an api key
in the header `x-api-key` \
To request for an API key, send an email to the author kamal2311@gmail.com

Issue a sample request

```shell script
curl --location --request GET 'https://zdn18ao2al.execute-api.us-east-2.amazonaws.com/test/v1/url-info/bad-url-1' \
--header 'x-api-key: <<API-Key>>' \
--header 'Content-Type: application/json'
```

responds with
 ```json
{
    "is_safe": false,
    "score": 8,
    "source": "Sophos"
}
```

and 

```shell script
curl --location --request GET 'https://zdn18ao2al.execute-api.us-east-2.amazonaws.com/test/v1/url-info/good-url' \
--header 'x-api-key: <<API-Key>>' \
--header 'Content-Type: application/json'
```
will respond with

```json
{
  "is_safe": true
}
```
## Checkout code
```shell script
git clone git@github.com:kamal2311/url-checker-lambda.git
cd url-checker-lambda
```


## Run some tests
```shell script
# Set your AWS profile
export AWS_PROFILE=kamal
./scripts/run-tests.sh
``` 

## Update lambda function code
```shell script
# Set your AWS profile
export AWS_PROFILE=kamal
./scripts/redeploy.sh
``` 
## Design and architecture reasoning

- Ensuring low latency
    - For the storage tier, we leverage a low latency, and infinitely scalable key value store for this (Dynamo DB in this case)
    - For the front-end tier we use API Gateway for the following reasons
        - API Gateway is a scalable serverless front-end that helps with concurrent requests with horizontal scalability
        - Takes care of SSL termination for us and our application code does not need to be burdened with extra CPU cycles spent in SSL connection.
        - We authorize requests with an API key that helps with rate limiting and usage analytics 
    
- Handling large number of requests
    - API gateway will scale up and down to accommodate for increase and decrease in traffic, keeping our infrastructure optimally utilized
    - This solution being serverless relieves us of operational burden which is of very much importanct at scale
- Handling a large list of urls in the malicious urls database
    - Dynamo db is ideally suited for meeting the low-latency response at large scale.
    - We have a very simple schema and a query pattern which is a good match for dynamo db
    
- Caching          
   - We can enable caching at the API gateway layer itself further helping with lowering the latency
   - We can also implement an in-memory cache with a hashmap and a queue to avoid making repeated requests to the database. 
    NOTE: for our application though, dynamo-db responses were so fast that caching made a little difference to the response time.
    Our requests currently respond within ~400ms on the average, most of which is spent in SSL handshake and network connection.
   - Time spent in sever processing is already in single/double digit milliseconds.
    
- Choice of Golang for the lambda function
   - Golang runtime has a fast cold start
   - It has a small memory footprint which helps with lowering the cost of lambda usage
    
- Supporting URL ingestion
   - We expect ~5000 new urls a day arriving every 10 minutes,
   - Assuming each URL and metadata on the average is of size 500 bytes, we will need to handle a POST body load of (5000 / 24 * 6 ) = ~35 Urls every 10 minutes
    This is a small payload 35 * 500 bytes = 17.5 KB for every POST request and can be easily handled.
   - We will generate a SHA-1 hash of the url to be stored in the db along with its metadata such as the source of the intelligence and maliciousness score 
   on a scale of 1 to 10.
     
## AWS infrastructure deployment components
- AWS API Gateway
- AWS Lambda proxy integration
- AWS Lambda function
- AWS Lambda execution role
- AWS Dynamo DB table
- AWS Cloudwatch logs and metrics

For this project , all these components were deployed using AWS console.
Future iterations may leverage AWS SAM CLI or Cloud formation or Terraform.

Some level of automation has been implemented through the following shell scripts.
The following script will generate a lambda function and a dynamodb table

```shell script
# Set your AWS profile
export AWS_PROFILE=kamal
./scripts/deploy.sh
```




