package encoding

import "github.com/pelletier/go-toml"

type tomlCodec struct{}

func (tomlCodec) Encode(v interface{}) ([]byte, error) {
	if m, ok := v.(map[string]interface{}); ok {
		t, err := toml.TreeFromMap(m)
		if err != nil {
			return nil, err
		}

		s, err := t.ToTomlString()
		if err != nil {
			return nil, err
		}

		return []byte(s), nil
	}

	return toml.Marshal(v)
}

func (tomlCodec) Decode(b []byte, v interface{}) error {
	tree, err := toml.LoadBytes(b)
	if err != nil {
		return err
	}

	if m, ok := v.(*map[string]interface{}); ok {
		vmap := *m
		tmap := tree.ToMap()
		for k, v := range tmap {
			vmap[k] = v
		}

		return nil
	}

	return tree.Unmarshal(v)
}
