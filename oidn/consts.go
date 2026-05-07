package oidn

// FilterType 滤波器类型（用于 NewFilter）
const (
	FilterTypeRT         = "RT"
	FilterTypeRTLightmap = "RTLightmap"
)

// Image param names (用于 SetImage)
const (
	ImageColor       = "color"  // 输入：带噪声的渲染图
	ImageOutput      = "output" // 输出：降噪结果
	ImageAlbedo      = "albedo" // 辅助：漫反射颜色
	ImageNormal      = "normal" // 辅助：世界空间法线（范围 -1..1）
	ImageDirectional = "directional"
	ImageLights      = "lights"
)

// Bool params (用于 SetBool)
const (
	BoolHDR      = "hdr"      // 输入图像是否为高动态范围
	BoolSRGB     = "srgb"     // 输入/输出是否使用 sRGB (GPU 后端影响精度)
	BoolCleanAux = "cleanAux" // 是否对辅助图像（albedo/normal）也执行降噪
)

// Int params (用于 SetInt)
const (
	IntMaxMemoryMB = "maxMemoryMB" // 滤波器最大内存使用（MB）
	IntQuality     = "quality"     // 质量模式（对应 Quality 枚举）
)
