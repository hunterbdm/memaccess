package memaccess

import (
	"errors"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type process struct {
	ID     uint32
	Name   string
	Handle windows.Handle
}

func (p *process) getModule(module string) (uintptr, error) {
	var (
		me32     windows.ModuleEntry32
		snap     windows.Handle
		szModule string
		err      error
	)

	snap, err = windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, p.ID)
	if err != nil {
		return 0, err
	}

	me32.Size = uint32(unsafe.Sizeof(me32))

	for err := windows.Module32First(snap, &me32); err != nil; err = windows.Module32Next(snap, &me32) {
		szModule = syscall.UTF16ToString(me32.Module[:])

		if szModule == module {
			return me32.ModBaseAddr, nil
		}
	}

	return me32.ModBaseAddr, errors.New("module not found")
}

func (p *process) open() error {
	handle, err := windows.OpenProcess(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|0xffff, false, p.ID)

	if err != nil {
		return err
	}

	p.Handle = handle

	return err
}

func getOpenProcessFromName(name string) (*process, error) {
	process, err := getProcessFromName(name)

	if err != nil {
		return nil, err
	}

	err = process.open()

	if err != nil {
		return nil, err
	}

	return process, nil
}

func getProcessFromName(name string) (*process, error) {
	pid, err := getProcessID(name)

	if err != nil {
		return nil, err
	}

	process := process{ID: pid, Name: name}

	return &process, nil
}

func getProcessID(process string) (uint32, error) {
	var (
		handle    windows.Handle
		pe32      windows.ProcessEntry32
		szExeFile string
	)

	pe32.Size = uint32(unsafe.Sizeof(pe32))

	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPALL, 0)
	if err != nil {
		return 0, err
	}

	for err := windows.Process32First(handle, &pe32); err == nil; err = windows.Process32Next(handle, &pe32) {
		szExeFile = syscall.UTF16ToString(pe32.ExeFile[:])

		if szExeFile == process {
			return pe32.ProcessID, nil
		}
	}

	return 0, errors.New("pid not found")
}
