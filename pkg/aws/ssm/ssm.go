// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package ssm

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// ParameterStore ARN's Resource includes the prefix 'parameter/'
func eliminatePrefix(resource string) string {
	if strings.HasPrefix(resource, "parameter/") {
		return resource[9:] // remove only 'parameter'
	}
	return resource
}

// GetParameter _
func GetParameter(region, name string) (*ssm.GetParameterOutput, error) {
	name = eliminatePrefix(name)
	svc := ssm.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	}
	return svc.GetParameter(input)
}

// GetParameterString _
func GetParameterString(region, name string) (string, error) {
	output, err := GetParameter(region, name)
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}

// PutParameter _
func PutParameter(region, name, value string) (*ssm.PutParameterOutput, error) {
	name = eliminatePrefix(name)
	svc := ssm.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &ssm.PutParameterInput{
		Type:      aws.String(ssm.ParameterTypeString),
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
	}
	return svc.PutParameter(input)
}
