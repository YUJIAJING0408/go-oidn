// Package oidn provides Go bindings for Intel Open Image Denoise.
//
// This package wraps the C API of Intel Open Image Denoise (OIDN), a library
// for filtering noise from images rendered with Monte Carlo ray tracing methods.
//
// Basic usage:
//
//	oidn.Init()
//	defer oidn.Shutdown()
//	device, _ := oidn.NewDevice(oidn.DeviceTypeCPU)
//	defer device.Release()
//	device.Commit()
//	// ... create buffers, filter, execute
package oidn
