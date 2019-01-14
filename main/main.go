package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kr/pretty"
)

// ParseClusterName parses a cluster name from the cluster ARN
func ParseClusterName(clusterARN string) (string, error) {
	matcher, err := regexp.Compile(":cluster/(.+)")
	if err != nil {
		log.Panic(err)
	}

	results := matcher.FindAllSubmatch([]byte(clusterARN), 1)
	if len(results) == 0 {
		return "", fmt.Errorf("could not parse cluster name from clusterARN '%s'", clusterARN)
	}
	return string(results[0][1]), nil
}

func handler(ctx context.Context, e events.CloudWatchEvent) error {
	pretty.Log(ctx)
	pretty.Log(e)
	var detail interface{}
	err := json.Unmarshal(e.Detail, &detail)
	if err != nil {
		return err
	}
	pretty.Log(detail)
	return nil
}

func main() {
	lambda.Start(handler)
}
