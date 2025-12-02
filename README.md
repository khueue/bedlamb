<p align="center">
  <img src="img/logo.png" alt="Bedlamb logo" width="200"/>
</p>

**TL;DR:** Kinda like curl, but for AWS Lambdas.

**NOTE: WORK IN PROGRESS.**

## Why?

I've recently been enamoured with running regular servers inside AWS Lambda:

- Develop your software as a regular server, listening to a port.
- Create a Dockerfile for it.
- Add a splash of AWS Lambda Web Adapter (literally one line of code):
  - It's great, I mean, look at it: https://github.com/awslabs/aws-lambda-web-adapter
- Deploy it.

Now, running it as a server has all the usual benefits of actually accessing it
through HTTP, locally and through some sort of ingress (Lambda URL etc.).
Running it locally as a server is pretty great, but want if I want to invoke a
deployed Lambda in AWS? It still expects some kind of event (not HTTP), right?

Since we now use the Web Adapter, when deployed as a Lambda it actually expects
an *API Gateway* event. So:

- **Bedlamb simply wraps your HTTP ambitions in an API Gateway event before
invoking your lambda using the AWS SDK**.

Like this:

```
bin/bedlamb \
  --path /health \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

Once again: *why?*

I admit, it's not rocket science. But what you get is the ability to work with
your deployed Lambdas just like regular servers, except they don't even have an
ingress. Fewer ingresses, fewer problems. Access is restricted through IAM.

I might even make it behave like a drop-in replacement for curl.

## Features

- Invoke AWS Lambda functions with API Gateway proxy request format
- Support for custom HTTP methods, paths, headers, query parameters, and request body
- curl-like interface for ease of use
- Uses AWS SDK Go v2 with standard credential providers

## Installation


Clone it and build it (you need Go):

```bash
git clone https://github.com/khueue/bedlamb.git
cd bedlamb
make build
```

## Usage

NOTE: The compiled binary is placed in `bin/`.

```
bedlamb [options] <lambda-arn>
```

### Options

- `-X, --method <method>` - HTTP method (default: GET)
- `-p, --path <path>` - Request path (default: /)
- `-H, --headers <headers>` - Headers in format 'Key1:Value1,Key2:Value2'
- `-d, --data <data>` - Request body data
- `-q, --query <query>` - Query string parameters in format 'key1=value1,key2=value2'
- `-v, --verbose` - Verbose output
- `--version` - Show version and exit

### Examples

**Simple GET request:**
```bash
bedlamb arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**POST request with JSON body:**
```bash
bedlamb \
  --method POST \
  --path /api/users \
  --data '{"name":"John","email":"john@example.com"}' \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**GET request with query parameters:**
```bash
bedlamb \
  --path /api/users \
  --query "page=1,limit=10" \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**Request with custom headers:**
```bash
bedlamb \
  --method POST \
  --path /api/data \
  --headers "Content-Type:application/json,Authorization:Bearer token123" \
  --data '{"key":"value"}' \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

**Verbose output:**
```bash
bedlamb \
  --verbose \
  --method GET \
  --path /health \
  arn:aws:lambda:us-east-1:123456789012:function:my-function
```

## AWS Credentials

The tool uses the AWS SDK's default credential chain, so fix your credentials. You will get a permission denied error from IAM if don't have access to invoke it.
