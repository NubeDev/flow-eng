package parameter

import (
	"github.com/francoispqt/gojay"
)

type Container struct {
	Parameters []NamedReadWriter
}

func NewContainer(parameters ...NamedReadWriter) *Container {
	return &Container{Parameters: parameters}
}

func (c *Container) Size() int {
	return len(c.Parameters)
}

func (c *Container) MarshalJSONArray(enc *gojay.Encoder) {
	for _, parameter := range c.Parameters {
		enc.Object(parameter)
	}
}

func (c *Container) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var p Base
	if err := dec.Object(&p); err != nil {
		return err
	}
	c.Parameters = append(c.Parameters, &p)
	return nil
}

func (c *Container) IsNil() bool {
	return len(c.Parameters) == 0
}
