package oidn

// DeviceType OIDN Device Type
type DeviceType int32

const (
	DeviceTypeDefault DeviceType = 0 // select device automatically
	DeviceTypeCPU     DeviceType = 1 // CPU device
	DeviceTypeSYCL    DeviceType = 2 // SYCL device
	DeviceTypeCUDA    DeviceType = 3 // CUDA device
	DeviceTypeHIP     DeviceType = 4 // HIP device
	DeviceTypeMetal   DeviceType = 5 // Metal device
)

// ErrorCode represents an OIDN error code.
type ErrorCode int32

const (
	ErrorNone             ErrorCode = 0
	ErrorUnknown          ErrorCode = 1
	ErrorInvalidArgument  ErrorCode = 2
	ErrorInvalidOperation ErrorCode = 3
	ErrorOutOfMemory      ErrorCode = 4
	ErrorUnsupportedHW    ErrorCode = 5
	ErrorCancelled        ErrorCode = 6
)

type Format int32

const (
	FormatUndefined Format = 0
	FormatFloat     Format = 1
	FormatFloat2    Format = 2
	FormatFloat3    Format = 3
	FormatFloat4    Format = 4
	FormatHalf      Format = 257
	FormatHalf2     Format = 258
	FormatHalf3     Format = 259
	FormatHalf4     Format = 260
)

// Quality represents filter quality/performance modes.
type Quality int32

const (
	QualityDefault  Quality = 0 // default quality
	QualityFast     Quality = 4 // high performance (interactive / real-time preview)
	QualityBalanced Quality = 5 // balanced quality/performance
	QualityHigh     Quality = 6 // high quality (final-frame rendering)
)
