package iso

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	parallelscommon "github.com/rickard-von-essen/packer-parallels/common"
	"path/filepath"
)

// This step creates the actual virtual machine.
//
// Produces:
//   vmName string - The name of the VM
type stepCreateVM struct {
	vmName string
}

func (s *stepCreateVM) Run(state multistep.StateBag) multistep.StepAction {

	config := state.Get("config").(*config)
	driver := state.Get("driver").(parallelscommon.Driver)
	ui := state.Get("ui").(packer.Ui)

	name := config.VMName
	path := filepath.Join(".", config.OutputDir)

	commands := make([][]string, 9)
	commands[0] = []string{
		"create", name,
		"--ostype", config.GuestOSType,
		"--distribution", config.GuestOSDistribution,
		"--dst", path,
		"--vmtype", "vm",
	}
	commands[1] = []string{"set", name, "--cpus", "1"}
	commands[2] = []string{"set", name, "--memsize", "512"}
	commands[3] = []string{"set", name, "--startup-view", "same"}
	commands[4] = []string{"set", name, "--on-shutdown", "close"}
	commands[5] = []string{"set", name, "--on-window-close", "keep-running"}
	commands[6] = []string{"set", name, "--auto-share-camera", "off"}
	commands[7] = []string{"set", name, "--device-del", "sound0"}
	commands[8] = []string{"set", name, "--smart-guard", "off"}

	ui.Say("Creating virtual machine...")
	for _, command := range commands {
		err := driver.Prlctl(command...)
		ui.Say(fmt.Sprintf("Doing: prlctl %s", command))
		if err != nil {
			err := fmt.Errorf("Error creating VM: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Set the VM name property on the first command
		if s.vmName == "" {
			s.vmName = name
		}
	}

	// Set the final name in the state bag so others can use it
	state.Put("vmName", s.vmName)
	return multistep.ActionContinue
}

func (s *stepCreateVM) Cleanup(state multistep.StateBag) {
	// TODO
	/*
	     if s.vmName == "" {
	   		return
	   	}

	   	driver := state.Get("driver").(parallelscommon.Driver)
	   	ui := state.Get("ui").(packer.Ui)
	   	config := state.Get("config").(*config)

	   	if config.DeleteVM {
	   		ui.Say("Deleting virtual machine...")
	   		if err := driver.Prlctl("delete", s.vmName); err != nil {
	   			ui.Error(fmt.Sprintf("Error deleting virtual machine: %s", err))
	   		}
	   	}
	*/
}
