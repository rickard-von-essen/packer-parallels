package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"log"
)

// Use defautl answers in Parallels
type StepUseDefaults struct{}

func (StepUseDefaults) Run(state multistep.StateBag) multistep.StepAction {
	log.Printf("DEFAULT: 1") //DEBUG
	driver := state.Get("driver").(Driver)
	log.Printf("DEFAULT: 2") //DEBUG
	ui := state.Get("ui").(packer.Ui)
	log.Printf("DEFAULT: 3")       //DEBUG
	log.Printf("State: %s", state) // DEBUG
	vmName := state.Get("vmName").(string)
	log.Printf("DEFAULT: 4") //DEBUG

	log.Println("Use default answers in Parallels")
	if err := driver.UseDefaults(vmName); err != nil {
		err := fmt.Errorf("Error configuring Parallels to use default answers: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (StepUseDefaults) Cleanup(multistep.StateBag) {}
