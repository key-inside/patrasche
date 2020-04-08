// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package secretsmanager

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretsManager ARN's Resource includes the prefix 'secret:'
func eliminatePrefix(resource string) string {
	parts := strings.Split(resource, ":")
	return parts[len(parts)-1]
}

// GetSecretValue _
func GetSecretValue(region, name string) (*secretsmanager.GetSecretValueOutput, error) {
	name = eliminatePrefix(name)
	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(name),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	return svc.GetSecretValue(input)
}

// GetSecretValueString _
func GetSecretValueString(region, name string) (string, error) {
	result, err := GetSecretValue(region, name)
	if err != nil {
		return "", err
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		return *result.SecretString, nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return "", err
	}
	return string(decodedBinarySecretBytes[:len]), nil
}

// PutSecretValue _
func PutSecretValue(region, name, value string) (*secretsmanager.PutSecretValueOutput, error) {
	name = eliminatePrefix(name)
	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(name),
		SecretString: aws.String(value),
	}
	return svc.PutSecretValue(input)
}
