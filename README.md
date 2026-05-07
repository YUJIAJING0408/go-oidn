# go-oidn

Go bindings for [Intel® Open Image Denoise](https://www.openimagedenoise.org/) (OIDN) — a high‑performance, easy‑to‑use library for denoising images rendered with ray tracing.

`go-oidn` 封装了 OIDN 的 C API，提供了符合 Go 语言习惯的设备、缓冲区、滤波器对象，并支持跨平台动态库加载。你可以将它作为库直接集成到自己的渲染器中，也可以使用内置的 CLI 工具快速对图片进行降噪。

## 项目状态

| 功能 | 状态          |
|------|-------------|
| 核心 API (设备 / 缓冲 / 滤波器) | ✅ 已完成       |
| Windows 动态库加载 | ✅ 已完成       |
| Linux / macOS 动态库加载 | ⏳ 计划中 |
| 图像 I/O (PNG 加载 / 保存) | ✅ 已完成       |
| 简单降噪示例 | ✅ 已完成       |
| CLI 命令行工具 (cobra) | 🚧 基础骨架已就绪  |
| 物理设备查询与枚举 | ⏳ 计划中       |
| CUDA / HIP / Metal 后端测试 | 🚧 计划中      |
| 单元测试与 CI | ⏳ 计划中       |

## 安装

首先确保你的系统已经安装了 Intel Open Image Denoise 运行时库（可从 [GITHUB](https://github.com/RenderKit/oidn/releases) 下载）。将动态库文件放置在项目约定的路径下，或通过环境变量 `OIDN_LIB_PATH` 指定目录。

```bash
go get github.com/YUJIAJING0408/go-oidn