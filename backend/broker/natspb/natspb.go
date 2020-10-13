package natspb

import (
	"errors"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type protoEncoder struct{}

var _ nats.Encoder = protoEncoder{}

func NewEncoder() nats.Encoder {
	return protoEncoder{}
}

func (p protoEncoder) Encode(_ string, v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("v is not proto.Message")
	}
	return proto.Marshal(msg)
}

func (p protoEncoder) Decode(_ string, data []byte, vPtr interface{}) error {
	msg, ok := vPtr.(proto.Message)
	if !ok {
		return errors.New("vPtr is not proto.Message")
	}
	return proto.Unmarshal(data, msg)
}
