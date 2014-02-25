package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"log"
	"path/filepath"
	"time"
)

// This step cleans up forwarded ports and exports the VM to an OVF.
//
// Uses:
//
// Produces:
//   exportPath string - The path to the resulting export.
type StepExport struct {
	OutputDir string
}

func (s *StepExport) Run(state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	// Wait a second to ensure VM is really shutdown.
	log.Println("1 second timeout to ensure VM is really shutdown")
	time.Sleep(1 * time.Second)

	// Clear out the Packer-created forwarding rule
	ui.Say("Preparing to export machine...")
	/*ui.Message(fmt.Sprintf(
		"Deleting forwarded port mapping for SSH (host port %d)",
		state.Get("sshHostPort")))
	command := []string{"modifyvm", vmName, "--natpf1", "delete", "packerssh"}
	if err := driver.Prlctl(command...); err != nil {
		err := fmt.Errorf("Error deleting port forwarding rule: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	*/

	outputPath := filepath.Join(s.OutputDir, vmName+".pvm")

	command := []string{
		"unregister",
		vmName,
	}

	ui.Say("Unregister virtual machine...")
	err := driver.Prlctl(command...)
	if err != nil {
		err := fmt.Errorf("Error unregistering virtual machine: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("exportPath", outputPath)

	return multistep.ActionContinue
}

func (s *StepExport) Cleanup(state multistep.StateBag) {}
