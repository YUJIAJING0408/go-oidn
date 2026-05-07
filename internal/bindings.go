package internal

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

type Functions struct {
	NewDevice              func(deviceType int32) uintptr
	CommitDevice           func(device uintptr)
	ReleaseDevice          func(device uintptr)
	SetDeviceErrorFunction func(device uintptr, cb uintptr, userPtr unsafe.Pointer)
	GetDeviceError         func(device uintptr, outMessage **byte) int32

	NewBuffer     func(device uintptr, byteSize uintptr) uintptr
	ReleaseBuffer func(buffer uintptr)
	WriteBuffer   func(buffer uintptr, byteOffset uintptr, byteSize uintptr, srcHostPtr unsafe.Pointer)
	ReadBuffer    func(buffer uintptr, byteOffset uintptr, byteSize uintptr, dstHostPtr unsafe.Pointer)

	NewFilter      func(device uintptr, filterType *byte) uintptr
	SetFilterImage func(filter uintptr, name *byte, buffer uintptr, format int32,
		width, height, byteOffset, pixelByteStride, rowByteStride uintptr)
	SetFilterBool  func(filter uintptr, name *byte, value bool)
	SetFilterInt   func(filter uintptr, name *byte, value int32)
	SetFilterFloat func(filter uintptr, name *byte, value float32)
	CommitFilter   func(filter uintptr)
	ExecuteFilter  func(filter uintptr)
	ReleaseFilter  func(filter uintptr)
	// Device
	GetNumPhysicalDevices   func() int32
	GetPhysicalDeviceInt    func(physicalDeviceID int32, name *byte) int32
	GetPhysicalDeviceString func(physicalDeviceID int32, name *byte) *byte
	GetPhysicalDeviceBool   func(physicalDeviceID int32, name *byte) bool
}

// F is the global instance holding all registered functions.
var F *Functions

// Init registers all OIDN C functions from the given library handle.
func Init(lib uintptr) error {
	f := &Functions{}
	purego.RegisterLibFunc(&f.NewDevice, lib, "oidnNewDevice")
	purego.RegisterLibFunc(&f.CommitDevice, lib, "oidnCommitDevice")
	purego.RegisterLibFunc(&f.ReleaseDevice, lib, "oidnReleaseDevice")
	purego.RegisterLibFunc(&f.SetDeviceErrorFunction, lib, "oidnSetDeviceErrorFunction")
	purego.RegisterLibFunc(&f.GetDeviceError, lib, "oidnGetDeviceError")
	purego.RegisterLibFunc(&f.NewBuffer, lib, "oidnNewBuffer")
	purego.RegisterLibFunc(&f.ReleaseBuffer, lib, "oidnReleaseBuffer")
	purego.RegisterLibFunc(&f.WriteBuffer, lib, "oidnWriteBuffer")
	purego.RegisterLibFunc(&f.ReadBuffer, lib, "oidnReadBuffer")
	purego.RegisterLibFunc(&f.NewFilter, lib, "oidnNewFilter")
	purego.RegisterLibFunc(&f.SetFilterImage, lib, "oidnSetFilterImage")
	purego.RegisterLibFunc(&f.SetFilterBool, lib, "oidnSetFilterBool")
	purego.RegisterLibFunc(&f.SetFilterInt, lib, "oidnSetFilterInt")
	purego.RegisterLibFunc(&f.SetFilterFloat, lib, "oidnSetFilterFloat")
	purego.RegisterLibFunc(&f.CommitFilter, lib, "oidnCommitFilter")
	purego.RegisterLibFunc(&f.ExecuteFilter, lib, "oidnExecuteFilter")
	purego.RegisterLibFunc(&f.ReleaseFilter, lib, "oidnReleaseFilter")
	purego.RegisterLibFunc(&f.GetNumPhysicalDevices, lib, "oidnGetNumPhysicalDevices")
	purego.RegisterLibFunc(&f.GetPhysicalDeviceInt, lib, "oidnGetPhysicalDeviceInt")
	purego.RegisterLibFunc(&f.GetPhysicalDeviceString, lib, "oidnGetPhysicalDeviceString")
	purego.RegisterLibFunc(&f.GetPhysicalDeviceBool, lib, "oidnGetPhysicalDeviceBool")
	F = f
	return nil
}
