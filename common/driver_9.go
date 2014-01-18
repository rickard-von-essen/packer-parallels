package common

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Parallels9Driver struct {
	// This is the path to the "prlctl" application.
	PrlctlPath string
}

func (d *Parallels9Driver) CreateSATAController(vmName string, name string) error {
	version, err := d.Version()
	if err != nil {
		return err
	}

	portCountArg := "--sataportcount"
	if strings.HasPrefix(version, "4.3") {
		portCountArg = "--portcount"
	}

	command := []string{
		"storagectl", vmName,
		"--name", name,
		"--add", "sata",
		portCountArg, "1",
	}
	command = []string{}
	return d.Prlctl(command...)
}

func (d *Parallels9Driver) Delete(name string) error {
	return d.Prlctl("delete", name)
}

func (d *Parallels9Driver) Import(name, path string) error {
	args := []string{
		"clone", path,
		"--regenerate-src-uuid",
	}

	return d.Prlctl(args...)
}

func (d *Parallels9Driver) IsRunning(name string) (bool, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(d.PrlctlPath, "list", name, "--no-header", "--output", "status")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return false, err
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		if line == "running" {
			return true, nil
		}

		if line == "suspended" {
			return true, nil
		}
		if line == "paused" {
			return true, nil
		}
	}

	return false, nil
}

func (d *Parallels9Driver) Stop(name string) error {
	if err := d.Prlctl("stop", name); err != nil {
		return err
	}

	// We sleep here for a little bit to let the session "unlock"
	time.Sleep(2 * time.Second)

	return nil
}

func (d *Parallels9Driver) UseDefaults(name string) error {
	if err := d.Prlctl("set", name, "--usedefanswers", "on"); err != nil {
		return err
	}

	return nil
}

func (d *Parallels9Driver) Prlctl(args ...string) error {
	var stdout, stderr bytes.Buffer

	log.Printf("Executing prlctl: %#v", args)
	cmd := exec.Command(d.PrlctlPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("prlctl error: %s", stderrString)
	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	return err
}

func (d *Parallels9Driver) Verify() error {
	log.Printf("Verifying Prlctl by doing nothing.")
	return nil
}

func (d *Parallels9Driver) Version() (string, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(d.PrlctlPath, "--version")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	versionOutput := strings.TrimSpace(stdout.String())
	log.Printf("prlctl --version output: %s", versionOutput)

	versionRe := regexp.MustCompile("[^.0-9]")
	matches := versionRe.Split(versionOutput, 2)
	if len(matches) == 0 || matches[0] == "" {
		return "", fmt.Errorf("No version found: %s", versionOutput)
	}

	log.Printf("prlctl version: %s", matches[0])
	return matches[0], nil
}
