package storage

import (
	"testing"
)

func TestAdapter(t *testing.T) {
	db := New("../flow.db")

	var a = NewAdapter(db)

	a.Add(&Connection{Name: "aaa"})

}
