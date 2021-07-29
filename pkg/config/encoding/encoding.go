package encoding

import (
	"errors"
)

type Decoder interface {
	Decode(b []byte, v interface{}) error
}

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

var (
	ErrEncoderNotFound = errors.New("encoder not found for config file format")
	ErrDecoderNotFound = errors.New("decoder not found for config file format")

	encoderRegistry = map[string]Encoder{}
	decoderRegistry = map[string]Decoder{}
)

func init() {
	{
		codec := yamlCodec{}
		encoderRegistry["yaml"] = codec
		encoderRegistry["yml"] = codec

		decoderRegistry["yaml"] = codec
		decoderRegistry["yml"] = codec
	}
	{
		codec := jsonCodec{}
		encoderRegistry["json"] = codec
		decoderRegistry["json"] = codec
	}
	{
		codec := tomlCodec{}
		encoderRegistry["toml"] = codec
		decoderRegistry["toml"] = codec
	}
}

// Encode encodes the value with specified format
func Encode(format string, v interface{}) ([]byte, error) {
	encoder, ok := encoderRegistry[format]
	if !ok {
		return nil, ErrEncoderNotFound
	}

	return encoder.Encode(v)
}

// Decode decodes the value with specified format
func Decode(format string, b []byte, v interface{}) error {
	decoder, ok := decoderRegistry[format]
	if !ok {
		return ErrDecoderNotFound
	}

	return decoder.Decode(b, v)
}
