package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/viper"
)

func DefaultConfig(optFns ...func(*config.LoadOptions) error) aws.Config {
	cfg, _ := config.LoadDefaultConfig(context.TODO(), optFns...)
	return cfg
}

func WithEndpoint(region, url string) config.LoadOptionsFunc {
	return config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: region,
			}, nil
		}))
}

func ext(resource string) string {
	for i := len(resource) - 1; i >= 0 && resource[i] != '/' && resource[i] != ':'; i-- {
		if resource[i] == '.' {
			return resource[i+1:]
		}
	}
	return ""
}

func GetConfigMap(arnStr string) (map[string]any, error) {
	str, err := GetString(arnStr)
	if err != nil {
		return nil, err
	}
	v := viper.New()
	v.SetConfigType(ext(arnStr))
	if err := v.ReadConfig(strings.NewReader(str)); err != nil {
		return nil, fmt.Errorf("failed to read from string reader: %w", err)
	}
	return v.AllSettings(), nil
}

func GetString(arnStr string) (string, error) {
	a, err := arn.Parse(arnStr)
	if err != nil {
		return "", err
	}
	switch a.Service {
	case "secretsmanager":
		return GetStringFromSecretsManager(DefaultConfig(), a)
	case "ssm":
		return GetStringFromParameterStore(DefaultConfig(), a)
	}
	return "", fmt.Errorf("not supported AWS service: %s", a.Service)
}

func GetStringFromSecretsManager(cfg aws.Config, a arn.ARN) (string, error) {
	client := secretsmanager.NewFromConfig(cfg)

	res, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(strings.TrimPrefix(a.Resource, "secret:")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret value: %w", err)
	}

	return *res.SecretString, nil
}

func GetStringFromParameterStore(cfg aws.Config, a arn.ARN) (string, error) {
	client := ssm.NewFromConfig(cfg)

	res, err := client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name: aws.String(strings.TrimPrefix(a.Resource, "parameter")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get parameter: %w", err)
	}

	return *res.Parameter.Value, nil
}

func GetItemFromDynamoDB(cfg aws.Config, tableName string, key map[string]any) (map[string]types.AttributeValue, error) {
	keyMap := map[string]types.AttributeValue{}
	var err error
	for k, v := range key {
		keyMap[k], err = attributevalue.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key: %w", err)
		}
	}

	client := dynamodb.NewFromConfig(cfg)
	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       keyMap,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return res.Item, nil
}

func PutItemToDynamoDB(cfg aws.Config, tableName string, item any) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      itemMap,
	})
	return err
}
