# movie-go

## Deploy
1. Build the binary `GOOS=linux go build -o main main.go`

2. Zip the binary
`zip deployment.zip main`

3. Create the lambda function
```
aws lambda create-function \
 --region us-east-2 \
 -- function-name DiscoverMovies \
 -- zip-file fileb://./deployment.zip \
-- runtime go1.x \
 --role arn:aws:iam::<account-id>:role/<role> \
 -- handler main
```