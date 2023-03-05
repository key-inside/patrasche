package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func Test_AWSParameterStore(t *testing.T) {
	arn := "arn:aws:ssm:::parameter/test/patrasche/config.yaml"
	str, err := GetString(arn)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else {
		t.Log(str)
	}
}

func Test_AWSSecretsManager(t *testing.T) {
	arn := "arn:aws:secretsmanager:::secret:test/patrasche/secrets.yaml"
	str, err := GetString(arn)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else {
		t.Log(str)
	}
}

func Test_AWSGetDynamoDBItem(t *testing.T) {
	cfg := DefaultConfig(WithEndpoint("eu-central-1", "http://localhost:8000"))
	key := map[string]any{"id": "dummy"}
	item, err := GetItemFromDynamoDB(cfg, ".meta", key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else {
		var blockNum uint64
		attributevalue.Unmarshal(item["block"], &blockNum)
		t.Log(blockNum)
	}
}

func Test_AWSPutDynamoDBItem(t *testing.T) {
	cfg := DefaultConfig(WithEndpoint("eu-central-1", "http://localhost:8000"))
	item := map[string]any{"id": "dummy", "block": 20190212}
	err := PutItemToDynamoDB(cfg, ".meta", item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
