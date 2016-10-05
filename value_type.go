package mruby

// ValueType is an enum of types that a Value can be and is returned by
// Value.Type().
type ValueType uint32

const (
	TypeFalse     ValueType = iota // 0
	TypeFree                       // 1
	TypeTrue                       // 2
	TypeFixnum                     // 3
	TypeSymbol                     // 4
	TypeUndef                      // 5
	TypeFloat                      // 6
	TypeCptr                       // 7
	TypeObject                     // 8
	TypeClass                      // 9
	TypeModule                     // 10
	TypeIClass                     // 11
	TypeSClass                     // 12
	TypeProc                       // 13
	TypeArray                      // 14
	TypeHash                       // 15
	TypeString                     // 16
	TypeRange                      // 17
	TypeException                  // 18
	TypeFile                       // 19
	TypeEnv                        // 20
	TypeData                       // 21
	TypeFiber                      // 22
	TypeMaxDefine                  // 23
	TypeNil       ValueType = 0xffffffff
)
