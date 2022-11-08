package node

func (n *Spec) GetMetadata() *Metadata {
	return n.Metadata
}

func (n *Spec) SetMetadata(m *Metadata) {
	n.Metadata = m
}

func (n *Spec) SetDynamicInputs() {
	if n.Metadata == nil {
		n.Metadata = &Metadata{}
	}
	n.Metadata.DynamicInputs = true
}

func (n *Spec) SetDynamicOutputs() {
	if n.Metadata == nil {
		n.Metadata = &Metadata{}
	}
	n.Metadata.DynamicOutputs = true
}
