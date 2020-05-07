// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package error

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResponsibleError is a responsible error from chaincodes
type ResponsibleError struct {
	Namespace string             `json:"namespace,omitempty"`
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	Causes    []ResponsibleError `json:"causes,omitempty"`
}

// Error implements error interface
func (e ResponsibleError) Error() string {
	return e.Message
}

// ToResponsibleError _
func ToResponsibleError(err error) (ResponsibleError, error) {
	msg := err.Error()
	idx := strings.LastIndex(msg, "Description: ")
	rErr := ResponsibleError{}
	if idx > -1 {
		idx += 13 // 13 = len("Description: ")
		err = json.Unmarshal([]byte(msg[idx:]), &rErr)
	} else {
		err = fmt.Errorf("not responsible error")
	}
	return rErr, err
}
