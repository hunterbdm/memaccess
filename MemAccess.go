package memaccess

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

/*
Credits:
github.com/gomem/gomem
golang.org/x/sys/windows
*/

type MemAccess struct {
	Proc *process
	Base uintptr
}

func NewMemAccess(process, module string) (*MemAccess, error) {
	proc, err := getOpenProcessFromName(process)
	if err != nil {
		return nil, err
	}

	baseAddress, err := proc.getModule(module)
	if err != nil {
		return nil, err
	}

	return &MemAccess{
		Proc: proc,
		Base: baseAddress,
	}, nil
}

// ReadPointer32 reads a 32-bit pointer from the process memory.
func (m *MemAccess) ReadPointer32(address uintptr) (uintptr, error) {
	v, err := m.ReadCustom(address, uint32(0))
	return uintptr(v.(uint32)), err
}

// ReadUInt32 reads a 32-bit unsigned integer from the process memory.
func (m *MemAccess) ReadUInt32(address uintptr) (uint32, error) {
	v, err := m.ReadCustom(address, uint32(0))
	return v.(uint32), err
}

// ReadByte reads a byte from the process memory.
func (m *MemAccess) ReadByte(address uintptr) (byte, error) {
	v, err := m.ReadCustom(address, byte(0))
	return v.(byte), err
}

// ReadPointerChain reads a pointer from a chain of offsets from the process memory.
// Starts at the module base address.
// It returns 0x0 if error occurs.
// Last offset is not dereference, as we assume it is a pointer to the value we want.
func (m *MemAccess) ReadPointerChain(chain ...uintptr) uintptr {
	lastPtr := m.Base

	for i := 0; i < len(chain)-1; i++ {
		newPtr, err := m.ReadPointer32(lastPtr + chain[i])
		if err != nil {
			return 0x0
		}

		lastPtr = newPtr
	}

	// Don't read pointer from last offset, as we assume this address points to the value we want.
	return lastPtr + chain[len(chain)-1]
}

// ReadCustom reads a custom value from the process memory.
func (m *MemAccess) ReadCustom(address uintptr, valueType any) (any, error) {
	var (
		value    any
		valuePtr = (*byte)(unsafe.Pointer(&value))
		size     = unsafe.Sizeof(valueType)
	)

	if err := windows.ReadProcessMemory(m.Proc.Handle, address, valuePtr, size, &size); err != nil {
		return nil, err
	}

	return value, nil
}

// Write writes a custom value to the process memory.
func (m *MemAccess) Write(offset uintptr, value any) error {
	var (
		valuePtr = (*byte)(unsafe.Pointer(&value))
		size     = unsafe.Sizeof(value)
	)

	if err := windows.WriteProcessMemory(m.Proc.Handle, offset, valuePtr, size, &size); err != nil {
		return err
	}

	return nil
}
