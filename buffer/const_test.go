package buffer

import "testing"

func BenchmarkConst_Write(b *testing.B) {
	b.ReportAllocs()

	buff := NewConst(Int8)
	toWrite := []byte{12}

	for i := 0; i < b.N; i++ {
		_, err := buff.Write(toWrite)
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkConst_Read(b *testing.B) {
	b.ReportAllocs()

	buff := NewConst(Int8)
	buff.Write([]byte{12})
	toRead := make([]byte, Int8)

	for i := 0; i < b.N; i++ {
		_, err := buff.Read(toRead)
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkConst_Copy(b *testing.B) {
	b.ReportAllocs()

	buffA := NewConst(Int8)
	buffA.Write([]byte{12})
	buffB := NewConst(Int8)

	for i := 0; i < b.N; i++ {
		_, err := buffA.Copy(buffB)
		if err != nil {
			b.Fail()
		}
	}
}
