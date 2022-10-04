package cmap

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var aa = map[string]interface{}{}
	a := New(aa)
	a.Add("a", 111)
	a.Add("b", 111)
	//a.Delete("b")
	fmt.Println(a.GetAll())
}
