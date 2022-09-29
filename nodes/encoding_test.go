package nodes

import (
	"encoding/json"
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"io/ioutil"
	"os"
	"testing"
)

// if there is no link then its not an array

func Test_test(t *testing.T) {

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
	//pprint.PrintJOSN(nodesParsed)

	decode, err := Decode(nodesParsed)

	pprint.PrintJOSN(decode)

}
