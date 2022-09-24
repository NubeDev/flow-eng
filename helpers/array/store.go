package array

import "fmt"

type ArrStore struct {
	s []interface{}
}

func (inst *ArrStore) IsEmpty() bool {
	length := len(inst.s)
	return length == 0
}

// Length of arr
func (inst *ArrStore) Length() int {
	length := len(inst.s)
	return length
}

// Print function
func (inst *ArrStore) Print() {
	length := len(inst.s)
	for i := 0; i < length; i++ {
		fmt.Print(inst.s[i], " ")
	}
	fmt.Println()
}

func (inst *ArrStore) Add(value interface{}) {
	inst.s = append(inst.s, value)
}

func (inst *ArrStore) RemoveLast() interface{} {
	length := len(inst.s)
	var res interface{}
	if length > 0 {
		res = inst.s[length-1]
		inst.s = inst.s[:length-1]
	}
	return res
}

func (inst *ArrStore) RemoveFirst() interface{} {
	length := len(inst.s)
	var res interface{}
	if length > 0 {
		res = inst.s[length-1]
		inst.s = inst.s[1:length]
	}
	return res
}

// Latest function
func (inst *ArrStore) Latest() interface{} {
	length := len(inst.s)
	res := inst.s[length-1]
	return res
}

func (inst *ArrStore) All() []interface{} {
	return inst.s
}

func (inst *ArrStore) GetByIndex(index int) interface{} {
	for i, v := range inst.s {
		if i == index {
			return v
		}
	}
	return nil
}
