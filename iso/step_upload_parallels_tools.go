package iso

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	parallelscommon "github.com/rickard-von-essen/packer-parallels/common"
	"log"
	"os"
	"strings"
)

type toolsPathTemplate struct {
	Version string
}

// This step uploads the guest additions ISO to the VM.
type stepUploadParallelsTools struct{}

func (s *stepUploadParallelsTools) Run(state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packer.Communicator)
	config := state.Get("config").(*config)
	driver := state.Get("driver").(parallelscommon.Driver)
	ui := state.Get("ui").(packer.Ui)

	// If we're attaching then don't do this, since we attached.
	if config.ParallelsToolsMode != ParallelsToolsModeUpload {
		log.Println("Not uploading Parallels Tools since mode is not upload")
		return multistep.ActionContinue
	}

	// Get the Parallels Tools path since we're doing it
	toolsPath := config.ParallelsToolsPath

	// Warning for https://github.com/mitchellh/packer/issues/951
	if strings.Contains(toolsPath, " ") {
		log.Printf("Space(s) found in tools path. You will hit https://github.com/mitchellh/packer/issues/951")
	}

	version, err := driver.Version()
	if err != nil {
		state.Put("error", fmt.Errorf("Error reading version for Parallels Tools upload: %s", err))
		return multistep.ActionHalt
	}

	f, err := os.Open(toolsPath)
	if err != nil {
		state.Put("error", fmt.Errorf("Error opening Parallels Tools ISO: %s", err))
		return multistep.ActionHalt
	}

	tplData := &toolsPathTemplate{
		Version: version,
	}

	config.ParallelsToolsPath, err = config.tpl.Process(config.ParallelsToolsPath, tplData)
	if err != nil {
		err := fmt.Errorf("Error preparing Parallels Tools path: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Uploading Parallels Tools ISO...")
	if err := comm.Upload(config.ParallelsToolsPath, f); err != nil {
		state.Put("error", fmt.Errorf("Error uploading Parallels Tools: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepUploadParallelsTools) Cleanup(state multistep.StateBag) {}
