export AWS_PROFILE=kamal
GOOS=linux go build -o main www.github.com/kamal2311/url-checker-lambda
~/go/bin/build-lambda-zip.exe -output main.zip main
aws lambda update-function-code --function-name my-function --zip-file fileb://main.zip