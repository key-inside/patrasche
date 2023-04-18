package block

import (
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"math"

	"github.com/hyperledger/fabric-protos-go/common"
)

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

// GenerateHash returns the ASN.1 marshaled hash bytes
func GenerateHash(block *common.Block) ([]byte, error) {
	header := block.Header
	asn1Header := asn1Header{
		PreviousHash: header.PreviousHash,
		DataHash:     header.DataHash,
	}
	if header.Number > uint64(math.MaxInt64) {
		return nil, errors.New("golang does not currently support encoding uint64 to asn1")
	}

	asn1Header.Number = int64(header.Number)
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write(result) // ignore error
	return hasher.Sum(nil), nil
}
