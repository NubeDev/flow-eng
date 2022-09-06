package parameter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/francoispqt/gojay"
	"unsafe"
)

type SerializableReadWriter interface {
	buffer.TypedReadWriter
	gojay.MarshalerJSONObject
	gojay.UnmarshalerJSONObject
}

type NamedReadWriter interface {
	SerializableReadWriter
	GetName() string
	GetDescription() string
}

type Base struct {
	buffer.Sync
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        Unit   `json:"unit"`
	ReadOnly    bool   `json:"readOnly"`
	UserVisible bool   `json:"userVisible"`
}

func New(_type buffer.Type, name, description string, unit Unit, readOnly, userVisible bool) *Base {
	return &Base{*buffer.NewSync(_type), name, description, unit, readOnly, userVisible}
}

func (p *Base) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("name", p.Name)
	enc.StringKey("description", p.Description)
	enc.StringKey("unit", *(*string)(unsafe.Pointer(&p.Unit)))
	enc.BoolKey("readOnly", p.ReadOnly)
	enc.BoolKey("userVisible", p.UserVisible)
	enc.AddArrayKey("raw", &p.Sync)
}

func (p *Base) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "name":
		return dec.String(&p.Name)
	case "raw":
		return dec.Array(&p.Sync)
	}
	return nil
}

func (p *Base) NKeys() int {
	return 0
}

func (p *Base) IsNil() bool {
	return p == nil
}

func (p *Base) GetName() string {
	return p.Name
}

func (p *Base) GetDescription() string {
	return p.Description
}
