package common

// A driver is able to talk to Parallels and perform certain
// operations with it. Some of the operations on here may seem overly
// specific, but they were built specifically in mind to handle features
// of the Parallels builder for Packer, and to abstract differences in
// versions out of the builder steps, so sometimes the methods are
// extremely specific.
type DriverVm interface {
	// Connect to the VM console. Should be done before sending key events.
	DisplayConnect() error

	// Disconnect from the VM console
	DisplayDisconnect()

	// Send a key scan codes via the Parallels Virtualization SDK C API
	SendKeyScanCode(uint32) error
}
