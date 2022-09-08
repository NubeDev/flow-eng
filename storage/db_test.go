package storage

import (
	"fmt"
	"testing"
)

func TestInitializeBuntDB(t *testing.T) {
	db := New("test.db")
	add, err := db.AddSettings(&Settings{})
	fmt.Println(add, err)

}
