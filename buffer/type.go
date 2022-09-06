package buffer

type Type int

const (
	Bool    Type = 1
	Int8    Type = 1
	Int16   Type = 2
	Int32   Type = 4
	Int64   Type = 8
	Uint8   Type = Int8
	Uint16  Type = Int16
	Uint32  Type = Int32
	Uint64  Type = Int64
	Float32 Type = 4
	Float64 Type = 8
	Byte    Type = Int8
	String  Type = 9
	Any     Type = 10
)
