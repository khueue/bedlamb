package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/google/uuid"
)

// version is set via ldflags at build time
var version = "dev"

// APIGatewayProxyRequest represents an API Gateway proxy request event
// {"rawPath":"/health","requestContext":{"http":{"method":"GET"}},"isBase64Encoded":false}
type APIGatewayProxyRequest struct {
	HTTPMethod            string            `json:"httpMethod"`
	Path                  string            `json:"path"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	Headers               map[string]string `json:"headers,omitempty"`
	Body                  string            `json:"body,omitempty"`
	IsBase64Encoded       bool              `json:"isBase64Encoded"`
	RequestContext        RequestContext    `json:"requestContext"`
}

// RequestContext represents the request context in an API Gateway event
type RequestContext struct {
	RequestID  string `json:"requestId"`
	Stage      string `json:"stage"`
	HTTPMethod string `json:"httpMethod"`
	Path       string `json:"path"`
}

func main() {
	// Define CLI flags with both short and long forms
	method := pflag.StringP("method", "X", "GET", "HTTP method")
	path := pflag.StringP("path", "p", "/", "Request path")
	headers := pflag.StringP("headers", "H", "", "Headers in format 'Key1:Value1,Key2:Value2'")
	data := pflag.StringP("data", "d", "", "Request body data")
	query := pflag.StringP("query", "q", "", "Query string parameters in format 'key1=value1,key2=value2'")
	verbose := pflag.BoolP("verbose", "v", false, "Verbose output")
	showVersion := pflag.Bool("version", false, "Show version and exit")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <lambda-arn>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -X POST --path /api/users -d '{\"name\":\"John\"}' arn:aws:lambda:us-east-1:123456789012:function:my-function\n", os.Args[0])
	}

	pflag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Check if Lambda ARN is provided
	if pflag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: Lambda ARN is required\n\n")
		pflag.Usage()
		os.Exit(1)
	}

	lambdaARN := pflag.Arg(0)

	// Parse headers
	headerMap := make(map[string]string)
	if *headers != "" {
		for pair := range strings.SplitSeq(*headers, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	// Parse query parameters
	queryMap := make(map[string]string)
	if *query != "" {
		for pair := range strings.SplitSeq(*query, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(parts) == 2 {
				queryMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	// Create API Gateway proxy request
	requestID := fmt.Sprintf("bedlamb-%s-%s", version, uuid.New().String())
	request := APIGatewayProxyRequest{
		HTTPMethod:            strings.ToUpper(*method),
		Path:                  *path,
		Headers:               headerMap,
		QueryStringParameters: queryMap,
		Body:                  *data,
		IsBase64Encoded:       false,
		RequestContext: RequestContext{
			RequestID:  requestID,
			Stage:      "prod",
			HTTPMethod: strings.ToUpper(*method),
			Path:       *path,
		},
	}

	// Marshal the request to JSON
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling request: %v\n", err)
		os.Exit(1)
	}

	// Pretty print the JSON payload
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, payload, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Lambda ARN: %s\n", lambdaARN)
		fmt.Fprintf(os.Stderr, "Request payload:\n%s\n\n", prettyJSON.String())
		fmt.Fprintf(os.Stderr, "Invoking Lambda...\n\n")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading AWS config: %v\n", err)
		os.Exit(1)
	}

	client := lambda.NewFromConfig(cfg)

	result, err := client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &lambdaARN,
		InvocationType: types.InvocationType("RequestResponse"),
		Payload:        payload,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error invoking Lambda: %v\n", err)
		os.Exit(1)
	}

	if result.FunctionError != nil {
		fmt.Fprintf(os.Stderr, "Lambda function error: %s\n", *result.FunctionError)
		fmt.Fprintf(os.Stderr, "Response: %s\n", string(result.Payload))
		os.Exit(1)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Status code: %d\n", result.StatusCode)
		fmt.Fprintf(os.Stderr, "Response:\n")
	}

	var prettyResponse bytes.Buffer
	if err := json.Indent(&prettyResponse, result.Payload, "", "  "); err != nil {
		// If it's not valid JSON, just print as-is
		fmt.Println(string(result.Payload))
	} else {
		fmt.Println(prettyResponse.String())
	}
}
