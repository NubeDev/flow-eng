package filter

const (
	category = "filter"
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
	onlyTrue           = "only-true"
	onlyFalse          = "only-false"
	preventNull        = "prevent-null"
	preventEqualFloat  = "prevent-equal-float"
	preventEqualString = "prevent-equal-string"
	onlyEqualFloat     = "only-equal-float"
	onlyEqualString    = "only-equal-string"
	onlyBetween        = "only-between"
	onlyGreater        = "only-greater"
	onlyLower          = "only-lower"
	preventDuplicates  = "prevent-duplicates"
)
