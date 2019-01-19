package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/kr/pretty"
)

type domainInfo struct {
	hostedZoneID string
	domain       string
}

// ParseClusterName parses a cluster name from the cluster ARN
func parseClusterName(clusterARN string) (string, error) {
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

// GetCluster retrieves the ecs.Cluster object for a named cluster in a given region
func getCluster(svc *ecs.ECS, name string) (*ecs.Cluster, error) {
	result, err := svc.DescribeClusters(
		&ecs.DescribeClustersInput{
			Clusters: []*string{aws.String(name)},
		},
	)
	if err != nil {
		return nil, err
	}

	clusters := result.Clusters
	if len(clusters) == 0 {
		return nil, fmt.Errorf("No clusters found named %s", name)
	}

	return result.Clusters[0], nil
}

// DomainInfo gets info on the domain that should point to this cluster's services
func getDomainInfo(svc *ecs.ECS, cluster *ecs.Cluster) (domainInfo, error) {
	var target domainInfo
	arn := cluster.ClusterArn
	tags, err := svc.ListTagsForResource(
		&ecs.ListTagsForResourceInput{
			ResourceArn: arn,
		},
	)
	if err != nil {
		return target, err
	}

	for _, tag := range tags.Tags {
		k := *tag.Key
		if k == "hostedZoneId" {
			target.hostedZoneID = *tag.Value
		}
		if k == "domain" {
			target.domain = *tag.Value
		}
	}
	if target.domain == "" {
		return target, fmt.Errorf("Missing hostedZoneId and/or domain tag for cluster with ARN %s", *arn)
	}

	return target, nil
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

	conf := aws.Config{}
	sess := session.Must(session.NewSession(&conf))
	svc := ecs.New(sess)

	cluster, err := getCluster(svc, "test-cluster")
	if err != nil {
		log.Panic(err)
	}

	domainInfo, err := getDomainInfo(svc, cluster)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
