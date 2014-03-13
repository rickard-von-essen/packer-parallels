package iso

import (
	"errors"
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	parallelscommon "github.com/rickard-von-essen/packer-parallels/common"
	"log"
	"strings"
)

const BuilderId = "rickard-von-essen.parallels"

// These are the different valid mode values for "parallels_tools_mode" which
// determine how guest additions are delivered to the guest.
const (
	ParallelsToolsModeDisable string = "disable"
	ParallelsToolsModeAttach         = "attach"
	ParallelsToolsModeUpload         = "upload"
)

type Builder struct {
	config config
	runner multistep.Runner
}

type config struct {
	common.PackerConfig                 `mapstructure:",squash"`
	parallelscommon.ExportConfig        `mapstructure:",squash"`
	parallelscommon.FloppyConfig        `mapstructure:",squash"`
	parallelscommon.OutputConfig        `mapstructure:",squash"`
	parallelscommon.RunConfig           `mapstructure:",squash"`
	parallelscommon.ShutdownConfig      `mapstructure:",squash"`
	parallelscommon.SSHConfig           `mapstructure:",squash"`
	parallelscommon.PrlctlConfig        `mapstructure:",squash"`
	parallelscommon.PrlctlVersionConfig `mapstructure:",squash"`

	BootCommand          []string `mapstructure:"boot_command"`
	DiskSize             uint     `mapstructure:"disk_size"`
	ParallelsToolsMode   string   `mapstructure:"parallels_tools_mode"`
	ParallelsToolsPath   string   `mapstructure:"parallels_tools_path"`
	ParallelsToolsSHA256 string   `mapstructure:"parallels_tools_sha256"`
	GuestOSType          string   `mapstructure:"guest_os_type"`
	GuestOSDistribution  string   `mapstructure:"guest_os_distribution"`
	HardDriveInterface   string   `mapstructure:"hard_drive_interface"`
	HTTPDir              string   `mapstructure:"http_directory"`
	HTTPPortMin          uint     `mapstructure:"http_port_min"`
	HTTPPortMax          uint     `mapstructure:"http_port_max"`
	ISOChecksum          string   `mapstructure:"iso_checksum"`
	ISOChecksumType      string   `mapstructure:"iso_checksum_type"`
	ISOUrls              []string `mapstructure:"iso_urls"`
	VMName               string   `mapstructure:"vm_name"`
	DeleteVM             bool     `mapstructure:"delete_vm"`

	RawSingleISOUrl string `mapstructure:"iso_url"`

	tpl *packer.ConfigTemplate
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {

	md, err := common.DecodeConfig(&b.config, raws...)
	if err != nil {
		return nil, err
	}

	b.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	b.config.tpl.UserVars = b.config.PackerUserVars

	// Accumulate any errors and warnings
	errs := common.CheckUnusedConfig(md)
	errs = packer.MultiErrorAppend(errs, b.config.ExportConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.FloppyConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(
		errs, b.config.OutputConfig.Prepare(b.config.tpl, &b.config.PackerConfig)...)
	errs = packer.MultiErrorAppend(errs, b.config.RunConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.ShutdownConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.SSHConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.PrlctlConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.PrlctlVersionConfig.Prepare(b.config.tpl)...)
	warnings := make([]string, 0)

	if b.config.DiskSize == 0 {
		b.config.DiskSize = 40000
	}

	if b.config.ParallelsToolsMode == "" {
		b.config.ParallelsToolsMode = "upload"
	}

	if b.config.ParallelsToolsPath == "" {
		b.config.ParallelsToolsPath = "prl-tools-lin.iso"
	}

	if b.config.HardDriveInterface == "" {
		b.config.HardDriveInterface = "sata"
	}

	if b.config.GuestOSType == "" {
		b.config.GuestOSType = "linux"
	}

	if b.config.HTTPPortMin == 0 {
		b.config.HTTPPortMin = 8000
	}

	if b.config.HTTPPortMax == 0 {
		b.config.HTTPPortMax = 9000
	}

	if b.config.VMName == "" {
		b.config.VMName = fmt.Sprintf("packer-%s", b.config.PackerBuildName)
	}

	// Errors
	templates := map[string]*string{
		"guest_tools_mode":      &b.config.ParallelsToolsMode,
		"guest_tools_sha256":    &b.config.ParallelsToolsSHA256,
		"guest_os_type":         &b.config.GuestOSType,
		"guest_os_distribution": &b.config.GuestOSDistribution,
		"hard_drive_interface":  &b.config.HardDriveInterface,
		"http_directory":        &b.config.HTTPDir,
		"iso_checksum":          &b.config.ISOChecksum,
		"iso_checksum_type":     &b.config.ISOChecksumType,
		"iso_url":               &b.config.RawSingleISOUrl,
		"vm_name":               &b.config.VMName,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = b.config.tpl.Process(*ptr, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	for i, url := range b.config.ISOUrls {
		var err error
		b.config.ISOUrls[i], err = b.config.tpl.Process(url, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing iso_urls[%d]: %s", i, err))
		}
	}

	validates := map[string]*string{
		"guest_tools_path": &b.config.ParallelsToolsPath,
	}

	for n, ptr := range validates {
		if err := b.config.tpl.Validate(*ptr); err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error parsing %s: %s", n, err))
		}
	}

	for i, command := range b.config.BootCommand {
		if err := b.config.tpl.Validate(command); err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("Error processing boot_command[%d]: %s", i, err))
		}
	}

	if b.config.HardDriveInterface != "ide" && b.config.HardDriveInterface != "sata" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("hard_drive_interface can only be ide or sata"))
	}

	if b.config.HTTPPortMin > b.config.HTTPPortMax {
		errs = packer.MultiErrorAppend(
			errs, errors.New("http_port_min must be less than http_port_max"))
	}

	if b.config.ISOChecksumType == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("The iso_checksum_type must be specified."))
	} else {
		b.config.ISOChecksumType = strings.ToLower(b.config.ISOChecksumType)
		if b.config.ISOChecksumType != "none" {
			if b.config.ISOChecksum == "" {
				errs = packer.MultiErrorAppend(
					errs, errors.New("Due to large file sizes, an iso_checksum is required"))
			} else {
				b.config.ISOChecksum = strings.ToLower(b.config.ISOChecksum)
			}

			if h := common.HashForType(b.config.ISOChecksumType); h == nil {
				errs = packer.MultiErrorAppend(
					errs,
					fmt.Errorf("Unsupported checksum type: %s", b.config.ISOChecksumType))
			}
		}
	}

	if b.config.RawSingleISOUrl == "" && len(b.config.ISOUrls) == 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("One of iso_url or iso_urls must be specified."))
	} else if b.config.RawSingleISOUrl != "" && len(b.config.ISOUrls) > 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("Only one of iso_url or iso_urls may be specified."))
	} else if b.config.RawSingleISOUrl != "" {
		b.config.ISOUrls = []string{b.config.RawSingleISOUrl}
	}

	for i, url := range b.config.ISOUrls {
		b.config.ISOUrls[i], err = common.DownloadableURL(url)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Failed to parse iso_url %d: %s", i+1, err))
		}
	}

	validMode := false
	validModes := []string{
		ParallelsToolsModeDisable,
		ParallelsToolsModeAttach,
		ParallelsToolsModeUpload,
	}

	for _, mode := range validModes {
		if b.config.ParallelsToolsMode == mode {
			validMode = true
			break
		}
	}

	if !validMode {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("parallels_tools_mode is invalid. Must be one of: %v", validModes))
	}

	if b.config.ParallelsToolsSHA256 != "" {
		b.config.ParallelsToolsSHA256 = strings.ToLower(b.config.ParallelsToolsSHA256)
	}

	// Warnings
	if b.config.ISOChecksumType == "none" {
		warnings = append(warnings,
			"A checksum type of 'none' was specified. Since ISO files are so big,\n"+
				"a checksum is highly recommended.")
	}

	if b.config.ShutdownCommand == "" {
		warnings = append(warnings,
			"A shutdown_command was not specified. Without a shutdown command, Packer\n"+
				"will forcibly halt the virtual machine, which may result in data loss.")
	}

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Create the driver that we'll use to communicate with Parallels
	driver, err := parallelscommon.NewDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating Parallels driver: %s", err)
	}

	steps := []multistep.Step{
		&common.StepDownload{
			Checksum:     b.config.ISOChecksum,
			ChecksumType: b.config.ISOChecksumType,
			Description:  "ISO",
			ResultKey:    "iso_path",
			Url:          b.config.ISOUrls,
		},
		&parallelscommon.StepOutputDir{
			Force: b.config.PackerForce,
			Path:  b.config.OutputDir,
		},
		&common.StepCreateFloppy{
			Files: b.config.FloppyFiles,
		},
		new(stepHTTPServer),
		new(stepCreateVM),
		new(stepCreateDisk),
		new(stepAttachISO),
		new(stepAttachParallelsTools),
		new(parallelscommon.StepAttachFloppy),
		&parallelscommon.StepPrlctl{
			Commands: b.config.Prlctl,
			Tpl:      b.config.tpl,
		},
		&parallelscommon.StepRun{
			BootWait: b.config.BootWait,
			Headless: b.config.Headless, // TODO: migth work on Enterprise Ed.
		},
		new(stepTypeBootCommand),
		&common.StepConnectSSH{
			SSHAddress:     parallelscommon.SSHAddress,
			SSHConfig:      parallelscommon.SSHConfigFunc(b.config.SSHConfig),
			SSHWaitTimeout: b.config.SSHWaitTimeout,
		},
		//&parallelscommon.StepUploadVersion{
		//	Path: b.config.PrlctlVersionFile,
		//},
		//new(stepUploadParallelsTools),
		new(common.StepProvision),
		&parallelscommon.StepShutdown{
			Command: b.config.ShutdownCommand,
			Timeout: b.config.ShutdownTimeout,
		},
		//new(parallelscommon.StepRemoveDevices),
		/*&parallelscommon.StepExport{
			OutputDir: b.config.OutputDir,
		},*/
	}

	// Setup the state bag
	state := new(multistep.BasicStateBag)
	state.Put("cache", cache)
	state.Put("config", &b.config)
	state.Put("driver", driver)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Run
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	//return parallelscommon.NewArtifact("test-dir")
	return parallelscommon.NewArtifact(b.config.OutputDir)
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
