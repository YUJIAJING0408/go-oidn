package oidn

import (
	"fmt"
	"unsafe"

	"github.com/YUJIAJING0408/go-oidn/internal"
	"github.com/ebitengine/purego"
)

// ErrorCallback is the function signature for device error callbacks.
type ErrorCallback func(userPtr unsafe.Pointer, code int32, message *byte) uintptr

// Device represents an OIDN device with automatic reference counting.
type Device struct {
	h uintptr
}

// Init loads the OIDN dynamic library and registers all C functions.
func Init() error {
	lib, err := loadLibrary()
	if err != nil {
		return err
	}
	return internal.Init(lib)
}

// NewDevice creates a new device of the specified type.
func NewDevice(deviceType DeviceType) (*Device, error) {
	if internal.F == nil {
		return nil, ErrNotInitialized
	}
	h := internal.F.NewDevice(int32(deviceType))
	if h == 0 {
		return nil, fmt.Errorf("oidn: failed to create device")
	}
	return &Device{h: h}, nil
}

// Release releases the device. After this call the device must not be used.
func (d *Device) Release() {
	if d.h != 0 {
		internal.F.ReleaseDevice(d.h)
		d.h = 0
	}
}

// Commit commits all previous changes to the device.
// Must be called before first using the device (e.g. creating filters).
func (d *Device) Commit() error {
	internal.F.CommitDevice(d.h)
	return d.GetError()
}

// SetErrorFunction sets the error callback for this device.
func (d *Device) SetErrorFunction(cb ErrorCallback) {
	cbPtr := purego.NewCallback(cb)
	internal.F.SetDeviceErrorFunction(d.h, cbPtr, unsafe.Pointer(nil))
}

// GetError returns and clears the first unqueried error for this device.
func (d *Device) GetError() error {
	var msgPtr *byte
	code := internal.F.GetDeviceError(d.h, &msgPtr)
	if code == 0 {
		return nil
	}
	msg := goString(msgPtr)
	return fmt.Errorf("oidn error %d: %s", code, msg)
}

// Handle returns the underlying C device handle (for advanced use).
func (d *Device) Handle() uintptr {
	return d.h
}
