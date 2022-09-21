package points

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestNewSync(t *testing.T) {
	s := NewSync()
	uuid := "a123"
	s.AddSync(uuid, 33)
	s.AddSync(uuid, 222)
	s.AddSync(uuid, 3333)
	s.AddSync(uuid, 3333)
	s.AddSync(uuid, 3333)
	s.AddSync(uuid, 3333)
	list := s.GetPoints()

	fmt.Println(len(list))

	pprint.Print(list)

	p := s.GetByPoint(uuid)
	pprint.PrintJOSN(p)

}
