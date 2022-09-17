package edge28

import "github.com/NubeDev/flow-eng/helpers/float"

func ProcessOutput(value *float64, ioType string, isDo bool) (float64, error) {
	value = limitValue(ioType, value)
	var err error
	var wv float64
	if isDo {
		if value != nil {
			writeValue := float.NonNil(value)
			wv, err = convertDigital(writeValue, false)
		}
	} else {
		wv, err = getValueUO(ioType)
	}
	return wv, err

}
