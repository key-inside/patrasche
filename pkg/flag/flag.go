// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package flag

import (
	"github.com/spf13/pflag"
)

// NewSSMFlagSet _
func NewSSMFlagSet(name, usage string) *pflag.FlagSet {
	flags := pflag.NewFlagSet(name, pflag.ExitOnError)
	flags.String(name+".region", "", "SSM region")
	flags.String(name+".parameter", "", "SSM parameter")
	return flags
}

// NewARNFlagSet _
func NewARNFlagSet(name, value, usage string) *pflag.FlagSet {
	flags := NewSSMFlagSet(name, usage)
	flags.String(name, value, usage)
	return flags
}

// NewARNFlagSetP _
func NewARNFlagSetP(name, shorthand, value, usage string) *pflag.FlagSet {
	flags := NewSSMFlagSet(name, usage)
	flags.StringP(name, shorthand, value, usage)
	return flags
}
