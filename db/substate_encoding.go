package db

import (
	"fmt"

	pb "github.com/Fantom-foundation/Substate/protobuf"
	"github.com/Fantom-foundation/Substate/rlp"
	"github.com/Fantom-foundation/Substate/substate"
	"github.com/Fantom-foundation/Substate/types"
	"github.com/golang/protobuf/proto"
)

// SetSubstateEncoding sets the runtime encoding/decoding behavior of substateDB
// intended usage:
//
//		db := &substateDB{..}
//	     db, err := db.SetSubstateEncoding(<schema>) // set encoding
//	     db.GetSubstateDecoder() // returns configured encoding
func (db *substateDB) SetSubstateEncoding(schema string) (*substateDB, error) {
	encoding, err := newSubstateEncoding(schema, db.GetCode)
	if err != nil {
		return nil, fmt.Errorf("Failed to set decoder; %w", err)
	}

	db.encoding = encoding
	return db, nil
}

// GetDecoder returns the encoding in use
func (db *substateDB) GetSubstateEncoding() string {
	if db.encoding == nil {
		return ""
	}
	return db.encoding.schema
}

type substateEncoding struct {
	schema string
	decode decoderFunc
}

// decoderFunc aliases the common function used to decode substate
type decoderFunc func([]byte, uint64, int) (*substate.Substate, error)

// codeLookup aliases codehash->code lookup necessary to decode substate
type codeLookup = func(types.Hash) ([]byte, error)

// newSubstateDecoder returns requested SubstateDecoder
func newSubstateEncoding(encoding string, lookup codeLookup) (*substateEncoding, error) {
	switch encoding {

	case "default", "rlp":
		return &substateEncoding{
			schema: "rlp",
			decode: func(bytes []byte, block uint64, tx int) (*substate.Substate, error) {
				return decodeRlp(bytes, lookup, block, tx)
			},
		}, nil

	case "protobuf", "pb":
		return &substateEncoding{
			schema: "protobuf",
			decode: func(bytes []byte, block uint64, tx int) (*substate.Substate, error) {
				return decodeProtobuf(bytes, lookup, block, tx)
			},
		}, nil

	default:
		return nil, fmt.Errorf("Encoding not supported: %s", encoding)

	}
}

// decodeSubstate defensively defaults to "default" if nil
func (db *substateDB) decodeToSubstate(bytes []byte, block uint64, tx int) (*substate.Substate, error) {
	if db.encoding == nil {
		db.SetSubstateEncoding("default")
	}
	return db.encoding.decode(bytes, block, tx)
}

// decodeRlp decodes into substate the provided rlp-encoded bytecode
func decodeRlp(bytes []byte, lookup codeLookup, block uint64, tx int) (*substate.Substate, error) {
	rlpSubstate, err := rlp.Decode(bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot decode substate data from rlp block: %v, tx %v; %w", block, tx, err)
	}

	return rlpSubstate.ToSubstate(lookup, block, tx)
}

// decodeProtobuf decodes into substate the provided rlp-encoded bytecode
func decodeProtobuf(bytes []byte, lookup codeLookup, block uint64, tx int) (*substate.Substate, error) {
	pbSubstate := &pb.Substate{}
	if err := proto.Unmarshal(bytes, pbSubstate); err != nil {
		return nil, fmt.Errorf("cannot decode substate data from protobuf block: %v, tx %v; %w", block, tx, err)
	}

	return pbSubstate.Decode(lookup, block, tx)
}

