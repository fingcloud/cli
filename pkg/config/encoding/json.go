package encoding

import "encoding/json"

type jsonCodec struct{}

func (jsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Decode(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
