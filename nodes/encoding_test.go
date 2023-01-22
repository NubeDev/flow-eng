package nodes

import (
	"encoding/json"
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Decode(t *testing.T) {

	var nodesParsed *NodesList
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)
	decode, err := Decode(nodesParsed)

	pprint.PrintJSON(FilterNodes(decode, FilterIsChild, ""))

}
