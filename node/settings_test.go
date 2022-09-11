package node

import (
	"testing"
)

type topic struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"topic"`
	Min      int    `json:"minLength" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

func TestBaseNode_GetSetting(t *testing.T) {

}
