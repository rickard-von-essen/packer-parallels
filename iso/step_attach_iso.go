package iso

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	parallelscommon "github.com/rickard-von-essen/packer-parallels/common"
)

// This step attaches the ISO to the virtual machine.
//
// Uses:
//
// Produces:
type stepAttachISO struct {
	diskPath string
}

func (s *stepAttachISO) Run(state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(parallelscommon.Driver)
	isoPath := state.Get("iso_path").(string)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	// Attach the disk to the controller
	command := []string{
		"set", vmName,
		"--device-set", "cdrom0",
		"--image", isoPath,
		"--enable", "--connect",
	}
	if err := driver.Prlctl(command...); err != nil {
		err := fmt.Errorf("Error attaching ISO: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Track the path so that we can unregister it from Parallels later
	s.diskPath = isoPath

	// Set some state so we know to remove
	state.Put("attachedIso", true)

	return multistep.ActionContinue
}

func (s *stepAttachISO) Cleanup(state multistep.StateBag) {
	if s.diskPath == "" {
		return
	}

	driver := state.Get("driver").(parallelscommon.Driver)
	vmName := state.Get("vmName").(string)

	command := []string{
		"set", vmName,
		"--device-del", "cdrom0",
	}

	// Remove the ISO. Note that this will probably fail since
	// stepRemoveDevices does this as well. No big deal.
	driver.Prlctl(command...)
}
