package edge28lib

import (
	"encoding/json"
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/services/clients/edgerest"
	"io/ioutil"
	"os"
	"testing"
)

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

func TestGetUI(t *testing.T) {
	ui, err := testUI()
	if err != nil {
		fmt.Println(err)
		return
	}
	is, err := processInput(ui.Val.UI1.Val, "thermistor_10k_type_2")
	fmt.Println(is, err)

	pprint.PrintJOSN(ui)
}

func TestGetDI(t *testing.T) {

	ui, err := testDI()

	if err != nil {
		fmt.Println(err)
		return
	}
	is, err := processInput(ui.Val.DI1.Val, "digital")
	fmt.Println(is, err)

}
