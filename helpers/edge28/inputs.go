package edge28

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/edge28/edgerest"
	"io/ioutil"
	"os"
)

func ProcessInput(val float64, ioType string) (float64, error) {
	return getValueUI(ioType, val)
}

func testUI() (*edgerest.UI, error) {
	var parsed *edgerest.UI
	jsonFile, err := os.Open("ui.json")
	if err != nil {
		fmt.Println(err)
		return nil, err

	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &parsed)
	return parsed, err

}

func testDI() (*edgerest.DI, error) {

	var parsed *edgerest.DI
	jsonFile, err := os.Open("di.json")
	if err != nil {
		fmt.Println(err)
		return nil, err

	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &parsed)
	return parsed, err

}
