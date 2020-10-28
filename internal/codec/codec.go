package codec

import (
	"io"
)

type Codec interface {
	Encode(interface{}, io.ReadWriter) error
	Decode(io.Reader, interface{}) error
	ContentType() string
	Accept() string
}
