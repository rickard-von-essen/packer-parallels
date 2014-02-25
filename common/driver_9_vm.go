package common

import (
	"github.com/rickard-von-essen/goprlapi"
)

type Parallels9DriverVm struct {
	vm goprlapi.VirtualMachine
}

func (v *Parallels9DriverVm) DisplayConnect() error {
	return v.vm.DisplayConnect()
}

func (v *Parallels9DriverVm) DisplayDisconnect() {
	v.vm.DisplayDisconnect()
}

func (v *Parallels9DriverVm) SendKeyScanCode(scancode uint32) error {
	return v.vm.SendKeyScanCode(scancode, 0)
}
