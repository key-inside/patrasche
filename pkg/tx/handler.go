// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package tx

// Handler _
type Handler interface {
	Handle(*Tx) error
}
