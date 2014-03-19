package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"log"
	"path/filepath"
	"time"
)

// This step unregister the virtual machine with Parallels
//
// Uses:
//
// Produces:
//   exportPath string - The path to the resulting virtual machine.
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

	outputPath := filepath.Join(s.OutputDir, vmName+".pvm")
	command := []string{"unregister", vmName}
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
