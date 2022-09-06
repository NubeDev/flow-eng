package buffer

import (
	"errors"
	"github.com/francoispqt/gojay"
	"io"
)

var (
	ErrTypesMismatch    = errors.New("provided buffers types are different")
	ErrDataSizeMismatch = errors.New("provided data size is different than buffer size")
)

type TypedReadWriter interface {
	io.ReadWriter
	Type() Type
}

type Const struct {
	data  []byte
	_type Type
}

func NewConst(_type Type) *Const {
	return &Const{make([]byte, _type), _type}
}

func (b *Const) MarshalJSONArray(enc *gojay.Encoder) {
	for _, value := range b.data {
		enc.Uint8(value)
	}
}

func (b *Const) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var value byte
	if err := dec.Uint8(&value); err != nil {
		return err
	}
	b.data = append(b.data, value)
	b._type = Type(len(b.data))
	return nil
}

func (b *Const) IsNil() bool {
	return b == nil
}

func (b *Const) Type() Type {
	return b._type
}

func (b *Const) Write(data []byte) (int, error) {
	if len(data) != len(b.data) {
		return 0, ErrDataSizeMismatch
	}
	copied := fastcopy(b.data, data)
	return copied, nil
}

func (b *Const) Read(data []byte) (int, error) {
	if len(data) != len(b.data) {
		return 0, ErrDataSizeMismatch
	}
	copied := fastcopy(data, b.data)
	return copied, nil
}

func (b *Const) Copy(other *Const) (int, error) {
	if b._type != other._type {
		return 0, ErrTypesMismatch
	}
	copied := fastcopy(other.data, b.data)
	return copied, nil
}

func fastcopy(dst []byte, src []byte) int {
	dstLen := len(dst)
	for i := 0; i < dstLen; i++ {
		dst[i] = src[i]
	}
	return dstLen
}
