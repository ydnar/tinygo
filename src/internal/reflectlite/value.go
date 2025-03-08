package reflectlite

import "unsafe"

type ValueFlags uint8

// These flags are shared with the reflect package.
const (
	ValueFlagIndirect ValueFlags = 1 << iota
	ValueFlagExported
	ValueFlagEmbedRO
	ValueFlagStickyRO

	ValueFlagRO = ValueFlagEmbedRO | ValueFlagStickyRO
)

func (v ValueFlags) RO() ValueFlags {
	if v&ValueFlagRO != 0 {
		return ValueFlagStickyRO
	}
	return 0
}

type Value struct {
	typecode *rawType
	value    unsafe.Pointer
	flags    ValueFlags
}

//go:linkname composeInterface runtime.composeInterface
func composeInterface(unsafe.Pointer, unsafe.Pointer) interface{}

//go:linkname decomposeInterface runtime.decomposeInterface
func decomposeInterface(i interface{}) (unsafe.Pointer, unsafe.Pointer)

func ValueOf(i interface{}) Value {
	typecode, value := decomposeInterface(i)
	return Value{
		typecode: (*rawType)(typecode),
		value:    value,
		flags:    ValueFlagExported,
	}
}

func (v Value) isIndirect() bool {
	return v.flags&ValueFlagIndirect != 0
}

func (v Value) isRO() bool {
	return v.flags&(ValueFlagRO) != 0
}

func (v Value) checkAddressable() {
	if !v.isIndirect() {
		panic("reflect: value is not addressable")
	}
}

func (v Value) checkRO() {
	if v.isRO() {
		panic("reflect: value is not settable")
	}
}

func (v Value) pointer() unsafe.Pointer {
	if v.isIndirect() {
		return *(*unsafe.Pointer)(v.value)
	}
	return v.value
}

func valueIsNil(v Value) bool {
	switch v.Kind() {
	case Chan, Map, Ptr, UnsafePointer:
		return v.pointer() == nil
	case Func:
		if v.value == nil {
			return true
		}
		fn := (*funcHeader)(v.value)
		return fn.Code == nil
	case Slice:
		if v.value == nil {
			return true
		}
		slice := (*sliceHeader)(v.value)
		return slice.data == nil
	case Interface:
		val := *(*interface{})(v.value)
		return val == nil
	default:
		panic(&ValueError{Method: "IsNil", Kind: v.Kind()})
	}
}

func valueLen(v Value) int {
	switch v.typecode.Kind() {
	case Array:
		return int(v.typecode.arrayLen())
	case Chan:
		return chanlen(v.pointer())
	case Map:
		return maplen(v.pointer())
	case Slice:
		return int((*sliceHeader)(v.value).len)
	case String:
		return int((*stringHeader)(v.value).len)
	default:
		panic(&ValueError{Method: "Len", Kind: v.Kind()})
	}
}

func valueElem(v Value) Value {
	switch v.Kind() {
	case Ptr:
		ptr := v.pointer()
		if ptr == nil {
			return Value{}
		}
		// Don't copy RO flags
		flags := (v.flags & (ValueFlagIndirect | ValueFlagExported)) | ValueFlagIndirect
		return Value{
			typecode: typeElem(v.typecode),
			value:    ptr,
			flags:    flags,
		}
	case Interface:
		typecode, value := decomposeInterface(*(*interface{})(v.value))
		return Value{
			typecode: (*rawType)(typecode),
			value:    value,
			flags:    v.flags &^ ValueFlagIndirect,
		}
	default:
		panic(&ValueError{Method: "Elem", Kind: v.Kind()})
	}
}

func valueString(v Value) string {
	switch v.Kind() {
	case String:
		// A string value is always bigger than a pointer as it is made of a
		// pointer and a length.
		return *(*string)(v.value)
	default:
		// Special case because of the special treatment of .String() in Go.
		return "<" + v.typecode.String() + " Value>"
	}
}

func valueField(v Value, i int) Value {
	if v.Kind() != Struct {
		panic(&ValueError{Method: "Field", Kind: v.Kind()})
	}
	structField := typeRawField(v.typecode, i)

	// Copy flags but clear EmbedRO; we're not an embedded field anymore
	flags := v.flags & ^ValueFlagEmbedRO
	if structField.PkgPath != "" {
		// No PkgPath => not exported.
		// Clear exported flag even if the parent was exported.
		flags &^= ValueFlagExported

		// Update the RO flag
		if structField.Anonymous {
			// Embedded field
			flags |= ValueFlagEmbedRO
		} else {
			flags |= ValueFlagStickyRO
		}
	} else {
		// Parent field may not have been exported but we are
		flags |= ValueFlagExported
	}

	size := v.typecode.Size()
	fieldType := structField.Type
	fieldSize := fieldType.Size()
	if v.isIndirect() || fieldSize > unsafe.Sizeof(uintptr(0)) {
		// v.value was already a pointer to the value and it should stay that
		// way.
		return Value{
			flags:    flags,
			typecode: fieldType,
			value:    unsafe.Add(v.value, structField.Offset),
		}
	}

	// The fieldSize is smaller than uintptr, which means that the value will
	// have to be stored directly in the interface value.

	if fieldSize == 0 {
		// The struct field is zero sized.
		// This is a rare situation, but because it's undefined behavior
		// to shift the size of the value (zeroing the value), handle this
		// situation explicitly.
		return Value{
			flags:    flags,
			typecode: fieldType,
			value:    unsafe.Pointer(nil),
		}
	}

	if size > unsafe.Sizeof(uintptr(0)) {
		// The value was not stored in the interface before but will be
		// afterwards, so load the value (from the correct offset) and return
		// it.
		ptr := unsafe.Add(v.value, structField.Offset)
		value := unsafe.Pointer(loadValue(ptr, fieldSize))
		return Value{
			flags:    flags &^ ValueFlagIndirect,
			typecode: fieldType,
			value:    value,
		}
	}

	// The value was already stored directly in the interface and it still
	// is. Cut out the part of the value that we need.
	value := maskAndShift(uintptr(v.value), structField.Offset, fieldSize)
	return Value{
		flags:    flags,
		typecode: fieldType,
		value:    unsafe.Pointer(value),
	}
}

var uint8Type = TypeOf(uint8(0)).(*rawType)

func valueIndex(v Value, i int) Value {
	switch v.Kind() {
	case Slice:
		// Extract an element from the slice.
		slice := *(*sliceHeader)(v.value)
		if uint(i) >= uint(slice.len) {
			panic("reflect: slice index out of range")
		}
		flags := (v.flags & (ValueFlagExported | ValueFlagIndirect)) | ValueFlagIndirect | v.flags.RO()
		elem := Value{
			typecode: typeElem(v.typecode),
			flags:    flags,
		}
		elem.value = unsafe.Add(slice.data, elem.typecode.Size()*uintptr(i)) // pointer to new value
		return elem
	case String:
		// Extract a character from a string.
		// A string is never stored directly in the interface, but always as a
		// pointer to the string value.
		// Keeping ValueFlagExported if set, but don't set ValueFlagIndirect
		// otherwise CanSet will return true for string elements (which is bad,
		// strings are read-only).
		s := *(*stringHeader)(v.value)
		if uint(i) >= uint(s.len) {
			panic("reflect: string index out of range")
		}
		return Value{
			typecode: uint8Type,
			value:    unsafe.Pointer(uintptr(*(*uint8)(unsafe.Add(s.data, i)))),
			flags:    v.flags & ValueFlagExported,
		}
	case Array:
		// Extract an element from the array.
		elemType := typeElem(v.typecode)
		elemSize := elemType.Size()
		size := v.typecode.Size()
		if size == 0 {
			// The element size is 0 and/or the length of the array is 0.
			return Value{
				typecode: typeElem(v.typecode),
				flags:    v.flags,
			}
		}
		if elemSize > unsafe.Sizeof(uintptr(0)) {
			// The resulting value doesn't fit in a pointer so must be
			// indirect. Also, because size != 0 this implies that the array
			// length must be != 0, and thus that the total size is at least
			// elemSize.
			addr := unsafe.Add(v.value, elemSize*uintptr(i)) // pointer to new value
			return Value{
				typecode: typeElem(v.typecode),
				flags:    v.flags,
				value:    addr,
			}
		}

		if size > unsafe.Sizeof(uintptr(0)) || v.isIndirect() {
			// The element fits in a pointer, but the array is not stored in the pointer directly.
			// Load the value from the pointer.
			addr := unsafe.Add(v.value, elemSize*uintptr(i)) // pointer to new value
			value := addr
			if !v.isIndirect() {
				// Use a pointer to the value (don't load the value) if the
				// 'indirect' flag is set.
				value = unsafe.Pointer(loadValue(addr, elemSize))
			}
			return Value{
				typecode: typeElem(v.typecode),
				flags:    v.flags,
				value:    value,
			}
		}

		// The value fits in a pointer, so extract it with some shifting and
		// masking.
		offset := elemSize * uintptr(i)
		value := maskAndShift(uintptr(v.value), offset, elemSize)
		return Value{
			typecode: typeElem(v.typecode),
			flags:    v.flags,
			value:    unsafe.Pointer(value),
		}
	default:
		panic(&ValueError{Method: "Index", Kind: v.Kind()})
	}
}

func valueSet(v, x Value) {
	v.checkAddressable()
	v.checkRO()
	if !x.typecode.AssignableTo(v.typecode) {
		panic("reflect.Value.Set: value of type " + x.typecode.String() + " cannot be assigned to type " + v.typecode.String())
	}

	if v.typecode.Kind() == Interface && x.typecode.Kind() != Interface {
		// move the value of x back into the interface, if possible
		if x.isIndirect() && x.typecode.Size() <= unsafe.Sizeof(uintptr(0)) {
			x.value = unsafe.Pointer(loadValue(x.value, x.typecode.Size()))
		}

		intf := composeInterface(unsafe.Pointer(x.typecode), x.value)
		x = Value{
			typecode: v.typecode,
			value:    unsafe.Pointer(&intf),
		}
	}

	size := v.typecode.Size()
	if size <= unsafe.Sizeof(uintptr(0)) && !x.isIndirect() {
		storeValue(v.value, size, uintptr(x.value))
	} else {
		memcpy(v.value, x.value, size)
	}
}

func (v Value) Type() Type {
	return v.typecode
}

// RawType is used in package runtime.
func (v Value) RawType() *rawType {
	return v.typecode
}

func (v Value) IsNil() bool {
	return valueIsNil(v)
}

func (v Value) Elem() Value {
	return valueElem(v)
}

func (v Value) NumField() int {
	return valueNumField(v)
}

func (v Value) Field(i int) Value {
	return valueField(v, i)
}

func (v Value) Index(i int) Value {
	return valueIndex(v, i)
}

func (v Value) Set(x Value) {
	valueSet(v, x)
}

func (v Value) Kind() Kind {
	return v.typecode.Kind()
}

func (v Value) Len() int {
	return valueLen(v)
}

func (v Value) String() string {
	return valueString(v)
}

type funcHeader struct {
	Context unsafe.Pointer
	Code    unsafe.Pointer
}

type sliceHeader struct {
	data unsafe.Pointer
	len  uintptr
	cap  uintptr
}

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

//go:linkname memcpy runtime.memcpy
func memcpy(dst, src unsafe.Pointer, size uintptr)

//go:linkname maplen runtime.hashmapLen
func maplen(p unsafe.Pointer) int

//go:linkname chanlen runtime.chanLen
func chanlen(p unsafe.Pointer) int
