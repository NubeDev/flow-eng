package db

import (
	"fmt"
	"testing"
)

func TestInitializeBuntDB(t *testing.T) {
	db := New("../flow.db")
	add, err := db.AddSettings(&Settings{})
	fmt.Println(add, err)

}
