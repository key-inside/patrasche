package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// the implicit contract of all these unmarshalers is that they
// will return a non-nil pointer whenever the error is nil

func wrapUnmarshalErr(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", msg, err)
	}
	return nil
}

// UnmarshalBlock unmarshals bytes to a Block
func UnmarshalBlock(bytes []byte) (*common.Block, error) {
	block := &common.Block{}
	return block, wrapUnmarshalErr(proto.Unmarshal(bytes, block), "Block")
}

// UnmarshalChaincodeDeploymentSpec unmarshals bytes to a ChaincodeDeploymentSpec
func UnmarshalChaincodeDeploymentSpec(bytes []byte) (*peer.ChaincodeDeploymentSpec, error) {
	cds := &peer.ChaincodeDeploymentSpec{}
	return cds, wrapUnmarshalErr(proto.Unmarshal(bytes, cds), "ChaincodeDeploymentSpec")
}

// UnmarshalChaincodeInvocationSpec unmarshals bytes to a ChaincodeInvocationSpec
func UnmarshalChaincodeInvocationSpec(bytes []byte) (*peer.ChaincodeInvocationSpec, error) {
	cis := &peer.ChaincodeInvocationSpec{}
	return cis, wrapUnmarshalErr(proto.Unmarshal(bytes, cis), "ChaincodeInvocationSpec")
}

// UnmarshalPayload unmarshals bytes to a Payload
func UnmarshalPayload(bytes []byte) (*common.Payload, error) {
	payload := &common.Payload{}
	return payload, wrapUnmarshalErr(proto.Unmarshal(bytes, payload), "Payload")
}

// UnmarshalEnvelope unmarshals bytes to a Envelope
func UnmarshalEnvelope(bytes []byte) (*common.Envelope, error) {
	envelope := &common.Envelope{}
	return envelope, wrapUnmarshalErr(proto.Unmarshal(bytes, envelope), "Envelope")
}

// UnmarshalChannelHeader unmarshals bytes to a ChannelHeader
func UnmarshalChannelHeader(bytes []byte) (*common.ChannelHeader, error) {
	chdr := &common.ChannelHeader{}
	return chdr, wrapUnmarshalErr(proto.Unmarshal(bytes, chdr), "ChannelHeader")
}

// UnmarshalChaincodeID unmarshals bytes to a ChaincodeID
func UnmarshalChaincodeID(bytes []byte) (*peer.ChaincodeID, error) {
	ccid := &peer.ChaincodeID{}
	return ccid, wrapUnmarshalErr(proto.Unmarshal(bytes, ccid), "ChaincodeID")
}

// UnmarshalSignatureHeader unmarshals bytes to a SignatureHeader
func UnmarshalSignatureHeader(bytes []byte) (*common.SignatureHeader, error) {
	sh := &common.SignatureHeader{}
	return sh, wrapUnmarshalErr(proto.Unmarshal(bytes, sh), "SignatureHeader")
}

// UnmarshalSerializedIdentity unmarshals bytes to a SerializedIdentity
func UnmarshalSerializedIdentity(bytes []byte) (*msp.SerializedIdentity, error) {
	sid := &msp.SerializedIdentity{}
	return sid, wrapUnmarshalErr(proto.Unmarshal(bytes, sid), "SerializedIdentity")
}

// UnmarshalHeader unmarshals bytes to a Header
func UnmarshalHeader(bytes []byte) (*common.Header, error) {
	hdr := &common.Header{}
	return hdr, wrapUnmarshalErr(proto.Unmarshal(bytes, hdr), "Header")
}

// UnmarshalChaincodeHeaderExtension unmarshals bytes to a ChaincodeHeaderExtension
func UnmarshalChaincodeHeaderExtension(bytes []byte) (*peer.ChaincodeHeaderExtension, error) {
	che := &peer.ChaincodeHeaderExtension{}
	return che, wrapUnmarshalErr(proto.Unmarshal(bytes, che), "ChaincodeHeaderExtension")
}

// UnmarshalProposalResponse unmarshals bytes to a ProposalResponse
func UnmarshalProposalResponse(bytes []byte) (*peer.ProposalResponse, error) {
	pr := &peer.ProposalResponse{}
	return pr, wrapUnmarshalErr(proto.Unmarshal(bytes, pr), "ProposalResponse")
}

// UnmarshalChaincodeAction unmarshals bytes to a ChaincodeAction
func UnmarshalChaincodeAction(bytes []byte) (*peer.ChaincodeAction, error) {
	ca := &peer.ChaincodeAction{}
	return ca, wrapUnmarshalErr(proto.Unmarshal(bytes, ca), "ChaincodeAction")
}

// UnmarshalResponse unmarshals bytes to a Response
func UnmarshalResponse(bytes []byte) (*peer.Response, error) {
	response := &peer.Response{}
	return response, wrapUnmarshalErr(proto.Unmarshal(bytes, response), "Response")
}

// UnmarshalChaincodeEvents unmarshals bytes to a ChaincodeEvent
func UnmarshalChaincodeEvents(bytes []byte) (*peer.ChaincodeEvent, error) {
	ce := &peer.ChaincodeEvent{}
	return ce, wrapUnmarshalErr(proto.Unmarshal(bytes, ce), "ChaicnodeEvent")
}

// UnmarshalProposalResponsePayload unmarshals bytes to a ProposalResponsePayload
func UnmarshalProposalResponsePayload(bytes []byte) (*peer.ProposalResponsePayload, error) {
	prp := &peer.ProposalResponsePayload{}
	return prp, wrapUnmarshalErr(proto.Unmarshal(bytes, prp), "ProposalResponsePayload")
}

// UnmarshalProposal unmarshals bytes to a Proposal
func UnmarshalProposal(bytes []byte) (*peer.Proposal, error) {
	prop := &peer.Proposal{}
	return prop, wrapUnmarshalErr(proto.Unmarshal(bytes, prop), "Proposal")
}

// UnmarshalTransaction unmarshals bytes to a Transaction
func UnmarshalTransaction(bytes []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	return tx, wrapUnmarshalErr(proto.Unmarshal(bytes, tx), "Transaction")
}

// UnmarshalChaincodeActionPayload unmarshals bytes to a ChaincodeActionPayload
func UnmarshalChaincodeActionPayload(bytes []byte) (*peer.ChaincodeActionPayload, error) {
	cap := &peer.ChaincodeActionPayload{}
	return cap, wrapUnmarshalErr(proto.Unmarshal(bytes, cap), "ChaincodeActionPayload")
}

// UnmarshalChaincodeProposalPayload unmarshals bytes to a ChaincodeProposalPayload
func UnmarshalChaincodeProposalPayload(bytes []byte) (*peer.ChaincodeProposalPayload, error) {
	cpp := &peer.ChaincodeProposalPayload{}
	return cpp, wrapUnmarshalErr(proto.Unmarshal(bytes, cpp), "ChaincodeProposalPayload")
}

// UnmarshalPayloadOrPanic unmarshals bytes to a Payload structure or panics
// on error
func UnmarshalPayloadOrPanic(bytes []byte) *common.Payload {
	payload, err := UnmarshalPayload(bytes)
	if err != nil {
		panic(err)
	}
	return payload
}

// UnmarshalEnvelopeOrPanic unmarshals bytes to an Envelope structure or panics
// on error
func UnmarshalEnvelopeOrPanic(bytes []byte) *common.Envelope {
	envelope, err := UnmarshalEnvelope(bytes)
	if err != nil {
		panic(err)
	}
	return envelope
}

// UnmarshalBlockOrPanic unmarshals bytes to an Block or panics
// on error
func UnmarshalBlockOrPanic(bytes []byte) *common.Block {
	block, err := UnmarshalBlock(bytes)
	if err != nil {
		panic(err)
	}
	return block
}

// UnmarshalChannelHeaderOrPanic unmarshals bytes to a ChannelHeader or panics
// on error
func UnmarshalChannelHeaderOrPanic(bytes []byte) *common.ChannelHeader {
	chdr, err := UnmarshalChannelHeader(bytes)
	if err != nil {
		panic(err)
	}
	return chdr
}

// UnmarshalSignatureHeaderOrPanic unmarshals bytes to a SignatureHeader or panics
// on error
func UnmarshalSignatureHeaderOrPanic(bytes []byte) *common.SignatureHeader {
	sighdr, err := UnmarshalSignatureHeader(bytes)
	if err != nil {
		panic(err)
	}
	return sighdr
}
