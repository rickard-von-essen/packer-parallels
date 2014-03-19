package vagrant

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/packer/packer"
)

// These are the extensions of files that are unnecessary for the function
// of a Parallels virtual machine.
var UnnecessaryFileExtensions = []string{".log", ".backup", ".Backup"}

type ParallelsProvider struct{}

func (p *ParallelsProvider) KeepInputArtifact() bool {
	return false
}

func (p *ParallelsProvider) Process(ui packer.Ui, artifact packer.Artifact, dir string) (vagrantfile string, metadata map[string]interface{}, err error) {
	// Create the metadata
	metadata = map[string]interface{}{"provider": "parallels"}

	// Copy all of the original contents into the temporary directory
	for _, path := range artifact.Files() {
		// If the file isn't critical to the function of the
		// virtual machine, we get rid of it.
		// It's done by the builder, but we need one more time
		// because unregistering a vm creates config.pvs.backup again.
		unnecessary := false
		ext := filepath.Ext(path)
		for _, unnecessaryExt := range UnnecessaryFileExtensions {
			if unnecessaryExt == ext {
				unnecessary = true
				break
			}
		}
		if unnecessary {
			continue
		}

		var pvmPath string

		tmpPath := filepath.ToSlash(path)
		pathRe := regexp.MustCompile(`^(.+?)([^/]+\.pvm/.+?)$`)
		matches := pathRe.FindStringSubmatch(tmpPath)
		if matches != nil {
			pvmPath = filepath.FromSlash(matches[2])
		} else {
			continue // Just copy a pvm
		}
		dstPath := filepath.Join(dir, pvmPath)

		ui.Message(fmt.Sprintf("Copying: %s", path))
		if err = CopyContents(dstPath, path); err != nil {
			return
		}
	}

	return
}
