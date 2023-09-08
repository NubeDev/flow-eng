package filter

const (
	Category = "filter"
)

// NODEs will single in/out
// only-true
// only-false
// prevent-null
// prevent-duplicates

// NODEs will double in
// only-between
// only-lower
// prevent-equal

const (
	preventNull        = "prevent-null"
	preventEqualFloat  = "prevent-equal-float"
	preventEqualString = "prevent-equal-string"
	onlyBetween        = "only-between"
	onlyGreater        = "only-greater"
	onlyLower          = "only-lower"
	preventDuplicates  = "prevent-duplicates"
)
