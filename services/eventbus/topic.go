package eventbus

// bus topics ...
const (
	All           = ".*"
	BacnetAll     = "bacnet.*"
	BacnetPV      = "bacnet.pv"
	BacnetPri     = "bacnet.pri"
	RubixIOInputs = "rubix.inputs"
)

// BusTopics return all bus topics
func BusTopics() []string {
	return []string{
		All,
		BacnetAll,
		BacnetPV,
		BacnetPri,
		RubixIOInputs,
	}
}
