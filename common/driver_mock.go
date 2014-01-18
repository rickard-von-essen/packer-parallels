package common

import "sync"

type DriverMock struct {
	sync.Mutex

	CreateSATAControllerVM         string
	CreateSATAControllerController string
	CreateSATAControllerErr        error

	DeleteCalled bool
	DeleteName   string
	DeleteErr    error

	ImportCalled bool
	ImportName   string
	ImportPath   string
	ImportErr    error

	IsRunningName   string
	IsRunningReturn bool
	IsRunningErr    error

	StopName string
	StopErr  error

	PrlctlCalls [][]string
	PrlctlErrs  []error

	VerifyCalled bool
	VerifyErr    error

	VersionCalled bool
	VersionResult string
	VersionErr    error
}

func (d *DriverMock) CreateSATAController(vm string, controller string) error {
	d.CreateSATAControllerVM = vm
	d.CreateSATAControllerController = vm
	return d.CreateSATAControllerErr
}

func (d *DriverMock) Delete(name string) error {
	d.DeleteCalled = true
	d.DeleteName = name
	return d.DeleteErr
}

func (d *DriverMock) Import(name, path string) error {
	d.ImportCalled = true
	d.ImportName = name
	d.ImportPath = path
	return d.ImportErr
}

func (d *DriverMock) IsRunning(name string) (bool, error) {
	d.Lock()
	defer d.Unlock()

	d.IsRunningName = name
	return d.IsRunningReturn, d.IsRunningErr
}

func (d *DriverMock) Stop(name string) error {
	d.StopName = name
	return d.StopErr
}

func (d *DriverMock) Prlctl(args ...string) error {
	d.PrlctlCalls = append(d.PrlctlCalls, args)

	if len(d.PrlctlErrs) >= len(d.PrlctlCalls) {
		return d.PrlctlErrs[len(d.PrlctlCalls)-1]
	}
	return nil
}

func (d *DriverMock) Verify() error {
	d.VerifyCalled = true
	return d.VerifyErr
}

func (d *DriverMock) Version() (string, error) {
	d.VersionCalled = true
	return d.VersionResult, d.VersionErr
}
