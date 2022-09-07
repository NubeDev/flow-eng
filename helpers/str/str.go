package str

func New(value string) *string {
	b := value
	return &b
}
func NonNil(b *string) string {
	if b == nil {
		return ""
	}
	return *b
}

func IsNil(b *string) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}
