package reflect

import (
	"internal/reflectlite"
	"unsafe"
)

type ChanDir = reflectlite.ChanDir

const (
	RecvDir = reflectlite.RecvDir
	SendDir = reflectlite.SendDir
	BothDir = reflectlite.BothDir
)

type Kind = reflectlite.Kind

const (
	Invalid       Kind = reflectlite.Invalid
	Bool          Kind = reflectlite.Bool
	Int           Kind = reflectlite.Int
	Int8          Kind = reflectlite.Int8
	Int16         Kind = reflectlite.Int16
	Int32         Kind = reflectlite.Int32
	Int64         Kind = reflectlite.Int64
	Uint          Kind = reflectlite.Uint
	Uint8         Kind = reflectlite.Uint8
	Uint16        Kind = reflectlite.Uint16
	Uint32        Kind = reflectlite.Uint32
	Uint64        Kind = reflectlite.Uint64
	Uintptr       Kind = reflectlite.Uintptr
	Float32       Kind = reflectlite.Float32
	Float64       Kind = reflectlite.Float64
	Complex64     Kind = reflectlite.Complex64
	Complex128    Kind = reflectlite.Complex128
	Array         Kind = reflectlite.Array
	Chan          Kind = reflectlite.Chan
	Func          Kind = reflectlite.Func
	Interface     Kind = reflectlite.Interface
	Map           Kind = reflectlite.Map
	Pointer       Kind = reflectlite.Pointer
	Slice         Kind = reflectlite.Slice
	String        Kind = reflectlite.String
	Struct        Kind = reflectlite.Struct
	UnsafePointer Kind = reflectlite.UnsafePointer
)

const Ptr = reflectlite.Ptr

type SelectCase = reflectlite.SelectCase

type StructField = reflectlite.StructField

type Type = reflectlite.Type

type Value = reflectlite.Value

type ValueError = reflectlite.ValueError

// Aliases of top-level exported functions in package reflectlite,
// sorted by name.

func Append(v Value, x ...Value) Value {
	return reflectlite.Append(v, x...)
}

func AppendSlice(s, t Value) Value {
	return reflectlite.AppendSlice(s, t)
}

func ArrayOf(n int, t Type) Type {
	return reflectlite.ArrayOf(n, t)
}

func Copy(dst, src Value) int {
	return reflectlite.Copy(dst, src)
}

func DeepEqual(x, y interface{}) bool {
	return reflectlite.DeepEqual(x, y)
}

func FuncOf(in, out []Type, variadic bool) Type {
	return FuncOf(in, out, variadic)
}

func Indirect(v Value) Value {
	return reflectlite.Indirect(v)
}

func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
	return reflectlite.MakeFunc(typ, fn)
}

func MakeMap(typ Type) Value {
	return reflectlite.MakeMap(typ)
}

func MakeMapWithSize(typ Type, n int) Value {
	return reflectlite.MakeMapWithSize(typ, n)
}

func MakeSlice(typ Type, len, cap int) Value {
	return reflectlite.MakeSlice(typ, len, cap)
}

func MapOf(key, value Type) Type {
	return MapOf(key, value)
}

func New(typ Type) Value {
	return reflectlite.New(typ)
}

func NewAt(typ Type, p unsafe.Pointer) Value {
	return reflectlite.NewAt(typ, p)
}

func PointerTo(t Type) Type {
	return reflectlite.PointerTo(t)
}

func PtrTo(t Type) Type {
	return reflectlite.PtrTo(t)
}

func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool) {
	return reflectlite.Select(cases)
}

func SliceOf(t Type) Type {
	return reflectlite.SliceOf(t)
}

func StructOf(fields []StructField) Type {
	return StructOf(fields)
}

func Swapper(slice interface{}) func(i, j int) {
	return reflectlite.Swapper(slice)
}

func TypeFor[T any]() Type {
	return reflectlite.TypeFor[T]()
}

func TypeOf(i interface{}) Type {
	return reflectlite.TypeOf(i)
}

func ValueOf(i interface{}) Value {
	return reflectlite.ValueOf(i)
}

func VisibleFields(t Type) []StructField {
	return reflectlite.VisibleFields(t)
}

func Zero(typ Type) Value {
	return reflectlite.Zero(typ)
}
