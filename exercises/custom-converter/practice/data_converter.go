package temporalconverters

import (
	"github.com/golang/snappy"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

var DataConverter = NewDataConverter(converter.GetDefaultDataConverter())

// NewDataConverter creates a new data converter that wraps the given data
// converter with snappy compression.
func NewDataConverter(underlying converter.DataConverter) converter.DataConverter {
	return converter.NewCodecDataConverter(underlying, NewPayloadCodec())
}

func NewPayloadCodec() converter.PayloadCodec {
	return &Codec{}
}

// Codec implements converter.PayloadEncoder for snappy compression.
type Codec struct{}

// Encode implements converter.PayloadCodec.Encode.
func (e *Codec) Encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		// TODO Part A: Use the first variable from `p.Marshal()` below.
		origBytes, err := p.Marshal()
		if err != nil {
			return payloads, err
		}

		// TODO Part A: Compress the marshalled variable to `b` using snappy.
		b := snappy.Encode(nil, origBytes)
		result[i] = &commonpb.Payload{
			Metadata: map[string][]byte{converter.MetadataEncoding: []byte("binary/snappy")},
			// TODO Part A: Assign `b` to `Data:` in this payload.
			Data: b,
		}
	}

	return result, nil
}

// Decode implements converter.PayloadCodec.Decode.
func (*Codec) Decode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		// Only if it's our encoding
		if string(p.Metadata[converter.MetadataEncoding]) != "binary/snappy" {
			result[i] = p
			continue
		}
		// Uncompress
		b, err := snappy.Decode(nil, p.Data)
		if err != nil {
			return payloads, err
		}
		// Unmarshal proto
		result[i] = &commonpb.Payload{}
		err = result[i].Unmarshal(b)
		if err != nil {
			return payloads, err
		}
	}

	return result, nil
}
