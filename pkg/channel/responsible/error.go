// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package responsible

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error is a responsible error from chaincodes
type Error struct {
	Namespace string  `json:"namespace,omitempty"`
	Code      string  `json:"code"`
	Message   string  `json:"message"`
	Causes    []Error `json:"causes,omitempty"`
}

// Error implements error interface
func (e Error) Error() string {
	return e.Message
}

// ToError _
func ToError(err error) (Error, error) {
	msg := err.Error()
	idx := strings.LastIndex(msg, "Description: ")
	rErr := Error{}
	if idx > -1 {
		idx += 13 // 13 = len("Description: ")
		err = json.Unmarshal([]byte(msg[idx:]), &rErr)
	} else {
		err = fmt.Errorf("not responsible error")
	}
	return rErr, err
}
