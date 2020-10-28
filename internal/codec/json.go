package codec

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	JSONContentType = "application/json"
	JSONAccept      = "application/json"
)

type JSON struct{}

func (s *JSON) Encode(source interface{}, target io.ReadWriter) error {
	enc := json.NewEncoder(target)
	enc.SetEscapeHTML(false)
	err := enc.Encode(source)

	return err
}

func (s *JSON) Decode(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)

	if errors.Is(err, io.EOF) {
		err = nil
	}

	return err
}

func (s *JSON) ContentType() string {
	return JSONContentType
}

func (s *JSON) Accept() string {
	return JSONAccept
}
