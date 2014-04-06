package common

import (
	"github.com/mitchellh/multistep"
	"testing"
)

func TestStepExport_impl(t *testing.T) {
	var _ multistep.Step = new(StepExport)
}

func TestStepExport(t *testing.T) {
	state := testState(t)
	step := new(StepExport)

	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)

	// Test the run
	if action := step.Run(state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}
	if _, ok := state.GetOk("error"); ok {
		t.Fatal("should NOT have error")
	}

	// Test output state
	if _, ok := state.GetOk("exportPath"); !ok {
		t.Fatal("should set exportPath")
	}

	// Test driver
	if len(driver.PrlctlCalls) != 1 {
		t.Fatal("should call prlctl")
	}
	if driver.PrlctlCalls[0][0] != "unregister" {
		t.Fatal("bad")
	}
}
