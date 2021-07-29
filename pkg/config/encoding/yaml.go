package encoding

import "gopkg.in/yaml.v3"

type yamlCodec struct{}

func (yamlCodec) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (yamlCodec) Decode(b []byte, v interface{}) error {
	return yaml.Unmarshal(b, v)
}
