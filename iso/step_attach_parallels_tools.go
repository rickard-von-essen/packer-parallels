package iso

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	parallelscommon "github.com/rickard-von-essen/packer-parallels/common"
	"log"
)

// This step attaches the Parallels Tools as a inserted CD onto
// the virtual machine.
//
// Uses:
//   config *config
//   driver Driver
//   parallels_tools_path string
//   ui packer.Ui
//   vmName string
//
// Produces:
type stepAttachParallelsTools struct {
	toolsPath string
}

func (s *stepAttachParallelsTools) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*config)
	driver := state.Get("driver").(parallelscommon.Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	// If we're not attaching the guest additions then just return
	if config.ParallelsToolsMode != ParallelsToolsModeAttach {
		log.Println("Not attaching parallels tools since we're uploading.")
		return multistep.ActionContinue
	}

	// Get the parallels tools path since we're doing it
	toolsPath := state.Get("parallels_tools_path").(string)

	// Attach the guest additions to the computer
	log.Println("Attaching Parallels Tools ISO onto IDE controller...")
	command := []string{
		"set", vmName,
		"--device-add", "cdrom",
		"--position", "0:0",
		"--connect",
		"--image", toolsPath,
	}
	if err := driver.Prlctl(command...); err != nil {
		err := fmt.Errorf("Error attaching Parallels Tools: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Track the path so that we can unregister it from Parallels later
	s.toolsPath = toolsPath

	return multistep.ActionContinue
}

func (s *stepAttachParallelsTools) Cleanup(state multistep.StateBag) {
	if s.toolsPath == "" {
		return
	}

	driver := state.Get("driver").(parallelscommon.Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	command := []string{
		"set", vmName,
		"--device-set", "cdrom:0:0",
		"--position", "0:0",
		"--disconnect",
	}

	if err := driver.Prlctl(command...); err != nil {
		ui.Error(fmt.Sprintf("Error unregistering guest additions: %s", err))
	}
}
