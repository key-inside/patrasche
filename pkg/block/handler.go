// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package block

import "github.com/hyperledger/fabric-protos-go/common"

// Handler _
type Handler interface {
	HandleBlock(hash []byte, block *common.Block) error
}

// Keeper _
type Keeper interface {
	LoadBlockNumber() (uint64, error)
	SaveBlockNumber(uint64) error
}
