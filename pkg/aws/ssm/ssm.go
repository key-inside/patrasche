// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetParameter _
func GetParameter(region, name string) (*ssm.GetParameterOutput, error) {
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
	svc := ssm.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &ssm.PutParameterInput{
		Type:      aws.String(ssm.ParameterTypeString),
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
	}
	return svc.PutParameter(input)
}
