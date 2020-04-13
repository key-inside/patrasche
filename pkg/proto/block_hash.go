// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package proto

import (
	"encoding/asn1"
	"fmt"
	"math"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/util"
)

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

// GenerateBlockHash returns the ASN.1 marshaled hash bytes
func GenerateBlockHash(block *common.Block) ([]byte, error) {
	header := block.Header
	asn1Header := asn1Header{
		PreviousHash: header.PreviousHash,
		DataHash:     header.DataHash,
	}
	if header.Number > uint64(math.MaxInt64) {
		return nil, fmt.Errorf("Golang does not currently support encoding uint64 to asn1")
	}

	asn1Header.Number = int64(header.Number)
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		return nil, err
	}

	return util.ComputeSHA256(result), nil
}
