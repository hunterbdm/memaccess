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

// Read reads a custom value/size from the process memory.
func (m *MemAccess) Read(address uintptr, resultPtr unsafe.Pointer, size uintptr) error {
	return windows.ReadProcessMemory(m.Proc.Handle, address, (*byte)(resultPtr), size, &size)
}

// Write writes a custom value to the process memory.
func (m *MemAccess) Write(address uintptr, valuePtr unsafe.Pointer, size uintptr) error {
	return windows.WriteProcessMemory(m.Proc.Handle, address, (*byte)(valuePtr), size, &size)
}

// ReadByte reads a byte from the process memory.
func (m *MemAccess) ReadByte(address uintptr) (byte, error) {
	var result byte
	if err := m.Read(address, unsafe.Pointer(&result), unsafe.Sizeof(result)); err != nil {
		return 0, nil
	}
	return result, nil
}

// ReadUInt32 reads a 32-bit unsigned integer from the process memory.
func (m *MemAccess) ReadUInt32(address uintptr) (uint32, error) {
	var result uint32
	if err := m.Read(address, unsafe.Pointer(&result), unsafe.Sizeof(result)); err != nil {
		return 0, nil
	}
	return result, nil
}

// ReadPointer32 reads a 32-bit pointer from the process memory.
func (m *MemAccess) ReadPointer32(address uintptr) (uintptr, error) {
	v, err := m.ReadUInt32(address)
	return uintptr(v), err
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

// WriteByte writes a byte to the process memory.
func (m *MemAccess) WriteByte(address uintptr, value byte) error {
	return m.Write(address, unsafe.Pointer(&value), unsafe.Sizeof(value))
}

// WriteUInt32 writes a 32-bit unsigned integer to the process memory.
func (m *MemAccess) WriteUInt32(address uintptr, value uint32) error {
	return m.Write(address, unsafe.Pointer(&value), unsafe.Sizeof(value))
}
