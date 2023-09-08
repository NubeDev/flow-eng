package link

const (
	linkInputString  = "link-input-string"
	linkOutputString = "link-output-string"
	linkInputNum     = "link-input-number"
	linkOutputNum    = "link-output-number"
	linkInputBool    = "link-input-boolean"
	linkOutputBool   = "link-output-boolean"
	Category         = "link"
)

var db *Store

func getStore() *Store {
	if db == nil {
		db = &Store{}
	}
	return db
}
