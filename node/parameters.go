package node

func BuildParameters(parameters *Parameters) *Parameters {
	if parameters == nil {
		parameters = &Parameters{}
	}
	return parameters
}
