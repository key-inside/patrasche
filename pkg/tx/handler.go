// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package tx

// Handler _
type Handler interface {
	Handle(*Tx) error
}

// BlockKeeper _
type BlockKeeper interface {
	LoadBlockNumber() (uint64, error)
	SaveBlockNumber(uint64) error
}
