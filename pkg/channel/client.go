// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package channel

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

// Client _
type Client struct {
	*channel.Client
}

func newClient(ctx context.ChannelProvider, options ...channel.ClientOption) (*Client, error) {
	_client, err := channel.New(ctx, options...)
	if err != nil {
		return nil, err
	}
	return &Client{_client}, nil
}

// ISSUE: if the duplicated transaction was MVCC_CONFLICT ...
// // ExecuteOnce executes the tranasction only once
// // tx ID is composed with nonce + creator, so, nonce is unique only within same creator !
// func (cc *Client) ExecuteOnce(nonce []byte, request channel.Request) (channel.Response, error) {
// 	res, err := cc.ExecuteWithNonce(request, nonce)
// 	if err != nil {
// 		if isDuplicatedTxError(err) {
// 			return res, &duplicatedTransactionError{txID: string(res.TransactionID)}
// 		}
// 		return res, err
// 	}
// 	return res, nil
// }

// func isDuplicatedTxError(err error) bool {
// 	return strings.Index(err.Error(), "duplicate transaction found") > -1
// }

// // DuplicatedTransactionError _
// type DuplicatedTransactionError interface {
// 	Error() string
// 	TxID() string
// }

// type duplicatedTransactionError struct {
// 	txID string
// }

// func (e *duplicatedTransactionError) Error() string {
// 	return fmt.Sprintf("duplicate transaction found [%s]", e.txID)
// }

// func (e *duplicatedTransactionError) TxID() string {
// 	return e.txID
// }
