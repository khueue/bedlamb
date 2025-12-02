# bedlamb

A CLI tool for invoking AWS Lambda functions with API Gateway formatted events, designed for Lambdas using the AWS Lambda Web Adapter.

## Features

- Invoke AWS Lambda functions with API Gateway proxy request format
- Support for custom HTTP methods, paths, headers, query parameters, and request body
- curl-like interface for ease of use
- Uses AWS SDK Go v2 with standard credential providers

## Installation

```bash
go install github.com/khu/bedlamb@latest
```

Or build from source:

```bash
git clone https://github.com/khu/bedlamb.git
cd bedlamb
go build -o bedlamb
```

## Usage

```
bedlamb [options] <lambda-arn>
```

### Options

- `-X <method>` - HTTP method (default: GET)
- `-path <path>` - Request path (default: /)
- `-H <headers>` - Headers in format 'Key1:Value1,Key2:Value2'
- `-d <data>` - Request body data
- `-q <query>` - Query string parameters in format 'key1=value1,key2=value2'
- `-v` - Verbose output

### Examples

**Simple GET request:**
```bash
bedlamb arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**POST request with JSON body:**
```bash
bedlamb -X POST -path /api/users -d '{"name":"John","email":"john@example.com"}' \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**GET request with query parameters:**
```bash
bedlamb -path /api/users -q "page=1,limit=10" \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**Request with custom headers:**
```bash
bedlamb -X POST -path /api/data \
  -H "Content-Type:application/json,Authorization:Bearer token123" \
  -d '{"key":"value"}' \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**Verbose output:**
```bash
bedlamb -v -X GET -path /health \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

## AWS Credentials

The tool uses the AWS SDK's default credential chain, which checks:
1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role (if running on EC2/ECS/Lambda)

Make sure you have appropriate permissions to invoke Lambda functions:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "lambda:InvokeFunction",
      "Resource": "arn:aws:lambda:*:*:function:*"
    }
  ]
}
```

## API Gateway Event Format

The tool formats requests as API Gateway proxy events with the following structure:

```json
{
  "httpMethod": "GET",
  "path": "/",
  "headers": {},
  "queryStringParameters": {},
  "body": "",
  "isBase64Encoded": false,
  "requestContext": {
    "requestId": "bedlamb-request",
    "stage": "prod",
    "httpMethod": "GET",
    "path": "/"
  }
}
```

## Requirements

1. CLI tool written in Go.
2. Tool should accept an AWS Lambda ARN as argument. It then uses the AWS SDK
  to invoke it. It should assume that the Lambda uses the AWS Lambda Web Adapter,
  so it expects input formatted as an API Gateway.

## License

MIT
