package link

type values struct {
	topic string
	value interface{}
}

type Store struct {
	s []*values
}

func (inst *Store) GetAll() []*values {
	return inst.s
}

func (inst *Store) Get(topic string) (value interface{}, found bool) {
	for _, v := range inst.s {
		if v.topic == topic {
			return v.value, true
		}
	}
	return nil, false
}

func (inst *Store) Add(topic string, value interface{}) {
	var existing bool
	for _, v := range inst.s {
		if v.topic == topic {
			existing = true
			v.value = value
		}
	}
	if !existing {
		inst.s = append(inst.s, &values{topic, value})
	}
}
