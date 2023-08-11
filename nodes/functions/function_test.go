package functions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/functions/rules"
	"testing"
)

func Test_runFunc(t *testing.T) {

	script := `

	RQL.Result = RQL.TimeWithMS()

`

	eng := rules.NewRuleEngine()
	n := "Core"
	props := make(rules.PropertiesMap)
	props[n] = eng
	client := "RQL"

	rule := &rules.RQL{
		Name:   "name",
		Script: script,
		Enable: true,
	}
	props[client] = rule

	err := eng.AddRule(rule, props)
	if err != nil {
		return
	}

	res, err := eng.ExecuteByName(rule.Name, true)

	if err != nil {
		fmt.Println("run", err)
		return
	}

	fmt.Println(res.String())

}
