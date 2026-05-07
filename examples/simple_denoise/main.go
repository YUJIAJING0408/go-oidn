package main

import (
	"fmt"

	"github.com/YUJIAJING0408/go-oidn/oidn"
)

func main() {
	// 初始化库
	if err := oidn.Init(); err != nil {
		panic(err)
	}
	for i := 0; i < oidn.GetNumPhysicalDevices(); i++ {
		typ := oidn.GetPhysicalDeviceInt(i, "type")
		name := oidn.GetPhysicalDeviceString(i, "name")
		fmt.Printf("Device %d: type=%d, name=%s\n", i, typ, name)
	}
	// 创建设备
	device, err := oidn.NewDevice(oidn.DeviceTypeCPU)
	if err != nil {
		panic(err)
	}
	defer device.Release()

	if err := device.Commit(); err != nil {
		panic(err)
	}

	// 加载图片
	colorData, width, height, err := oidn.LoadPNG("noisy.png")
	if err != nil {
		panic(err)
	}

	// 创建缓冲区
	colorBuf, _ := device.NewBuffer(len(colorData) * 4)
	defer colorBuf.Release()
	colorBuf.Write(colorData)

	outputBuf, _ := device.NewBuffer(width * height * 3 * 4)
	defer outputBuf.Release()

	// 创建滤波器
	filter, err := device.NewFilter("RT")
	if err != nil {
		panic(err)
	}
	defer filter.Release()

	filter.SetImage("color", colorBuf, oidn.FormatFloat3, width, height)
	filter.SetImage("output", outputBuf, oidn.FormatFloat3, width, height)
	filter.SetBool("hdr", false)
	filter.Commit()

	fmt.Println("Executing...")
	filter.Execute()

	if err := device.GetError(); err != nil {
		panic(err)
	}

	// 保存结果
	result := outputBuf.Read(width * height * 3)
	oidn.SavePNG(result, width, height, "denoised.png")
	fmt.Println("Done.")
}
