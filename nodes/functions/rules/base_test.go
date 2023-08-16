package rules

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"testing"
)

func TestFilterTest(t *testing.T) {
	body := `{"in1":"{\"in1\":\"107.00\",\"in2\":\"2\",\"in3\":null,\"in4\":null,\"in5\":null,\"in6\":null,\"in7\":null,\"in8\":null,\"in9\":null,\"in10\":null}","in2":"33","in3":null,"in4":null,"in5":null,"in6":null,"in7":null,"in8":null,"in9":null,"in10":null}`
	// body = `{\"in1\":\"107.00\",\"in2\":\"2\",\"in3\":null,\"in4\":null,\"in5\":null,\"in6\":null,\"in7\":null,\"in8\":null,\"in9\":null,\"in10\":null}`

	// body := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	// value := gjson.Get(body, "in1")
	// // value2 := gjson.Get(value.String(), "")
	//
	// fmt.Println(value.String())

	// const json = `{"name":{"first":"Tom","last":"Hanks"},"age":61}`
	name := gojsonq.New().FromString(body).Find("in1")
	name2 := gojsonq.New().FromString(fmt.Sprint(name)).Find("in1")

	fmt.Println(name2)
	fmt.Println(name)

}

func TestNewRuleEngine(t *testing.T) {

	script := `
	let a = RQL.in1+10
	RQL.Result = a+111111
`

	rule := &RQL{
		UUID:              "",
		Name:              "test",
		LatestRunDate:     "",
		Script:            script,
		Schedule:          "",
		Enable:            true,
		ResultStorageSize: 0,
		Result:            nil,
	}
	arg := map[string]interface{}{"in1": 22.2, "in2": 23, "in3": 23, "Result": nil}
	props := PropertiesMap{
		"RQL": arg,
	}

	r := NewRuleEngine()

	err := r.AddRule(rule, props)

	if err != nil {
		fmt.Println("add", err)
		return
	}

	res, err := r.ExecuteByName(rule.Name, true)

	if err != nil {
		fmt.Println("run", err)
		return
	}

	fmt.Println(res.String())

}
