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
	path := filepath.Join(".", config.OutputDir) //, fmt.Sprintf("%s.pvm", name)

	commands := make([][]string, 3)
	commands[0] = []string{
		"create", name,
		"--ostype", config.GuestOSType,
		"--distribution", config.GuestOSDistribution,
		"--dst", path,
	}
	commands[1] = []string{"set", name, "--cpus", "1"}
	commands[2] = []string{"set", name, "--memsize", "512"}

	ui.Say("Creating virtual machine...")
	for _, command := range commands {
		err := driver.Prlctl(command...)
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
}
