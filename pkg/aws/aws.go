// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package aws

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/aws/secretsmanager"
	"github.com/key-inside/patrasche/pkg/aws/ssm"
)

// GetARN returns a ARN from default viper
func GetARN(name string) (arnObj arn.ARN, raw string, err error) {
	return GetARNFromViper(viper.GetViper(), name)
}

// GetARNFromViper _
func GetARNFromViper(v *viper.Viper, name string) (arnObj arn.ARN, raw string, err error) {
	region := name + ".region"
	resource := name + ".parameter"
	if v.IsSet(region) && v.IsSet(resource) {
		return arn.ARN{
			Service:  "ssm",
			Region:   v.GetString(region),
			Resource: v.GetString(resource),
		}, "", nil
	}
	raw = v.GetString(name)
	arnObj, err = arn.Parse(raw)
	return
}

// GetContentTypeOfARN _
func GetContentTypeOfARN(arnObj arn.ARN) string {
	ext := filepath.Ext(arnObj.Resource)
	if "" == ext {
		return "json" // default is json
	}
	return ext[1:] // remove .
}

// GetViperWithARN _
func GetViperWithARN(arnObj arn.ARN) (*viper.Viper, error) {
	in, typ, err := GetReaderWithARN(arnObj)
	if err != nil {
		return nil, err
	}
	subV := viper.New()
	subV.SetConfigType(typ)
	if err := subV.ReadConfig(in); err != nil {
		return nil, err
	}
	return subV, nil
}

// GetStringWithARN _
func GetStringWithARN(arnObj arn.ARN) (string, error) {
	switch arnObj.Service {
	case "ssm":
		return ssm.GetParameterString(arnObj.Region, arnObj.Resource)
	case "secretsmanager":
		return secretsmanager.GetSecretValueString(arnObj.Region, arnObj.Resource)
	}
	return "", fmt.Errorf("not supported AWS service")
}

// GetReaderWithARN _
func GetReaderWithARN(arnObj arn.ARN) (io.Reader, string, error) {
	v, err := GetStringWithARN(arnObj)
	if err != nil {
		return nil, "", err
	}
	return strings.NewReader(v), GetContentTypeOfARN(arnObj), nil
}

// PutStringWithARN _
func PutStringWithARN(arnObj arn.ARN, value string) error {
	switch arnObj.Service {
	case "ssm":
		_, err := ssm.PutParameter(arnObj.Region, arnObj.Resource, value)
		return err
	case "secretsmanager":
		_, err := secretsmanager.PutSecretValue(arnObj.Region, arnObj.Resource, value)
		return err
	}
	return fmt.Errorf("not supported AWS service")
}
