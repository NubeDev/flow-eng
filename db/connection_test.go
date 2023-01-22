package db

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestNew(t *testing.T) {

	db := New("../flow.db")
	connection, err := db.AddConnection(&Connection{
		Name:        "test",
		Application: "flow",
	})
	pprint.PrintJSON(connection)
	fmt.Println(err)

	connections, err := db.GetConnections()
	fmt.Println(connections, err)
	if err != nil {
		return
	}
}
