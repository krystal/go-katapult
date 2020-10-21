package katapult

import (
	"encoding/json"
	"io"
)

const (
	JSONCodecContentType = "application/json"
	JSONCodecAccept      = "application/json"
)

type JSONCodec struct{}

func (s *JSONCodec) Encode(source interface{}, target io.ReadWriter) error {
	enc := json.NewEncoder(target)
	enc.SetEscapeHTML(false)
	err := enc.Encode(source)

	return err
}

func (s *JSONCodec) Decode(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)

	// ignore EOF errors caused by empty response body
	if err == io.EOF {
		err = nil
	}

	return err
}

func (s *JSONCodec) ContentType() string {
	return JSONCodecContentType
}

func (s *JSONCodec) Accept() string {
	return JSONCodecAccept
}
