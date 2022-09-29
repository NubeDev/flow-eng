package edge28lib

var UOTypes = struct {
	DIGITAL string
	VOLTSDC string
}{
	DIGITAL: "digital",
	VOLTSDC: "voltage_dc",
}

var UITypes = struct {
	RAW             string
	DIGITAL         string
	VOLTSDC         string
	MILLIAMPS       string
	RESISTANCE      string
	THERMISTOR10KT2 string
}{
	RAW:             "raw",
	DIGITAL:         "digital",
	VOLTSDC:         "voltage_dc",
	MILLIAMPS:       "current",
	THERMISTOR10KT2: "thermistor_10k_type_2",
}
