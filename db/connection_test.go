package db

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestNew(t *testing.T) {

	db := New("../flow.db")
	connection, err := db.AddConnection(&Connection{
		Application: "mqtt",
	})
	pprint.PrintJOSN(connection)
	fmt.Println(err)

	connections, err := db.GetConnections()
	fmt.Println(connections, err)
	if err != nil {
		return
	}
}
