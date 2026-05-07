package multi_round_denoising

import (
	"fmt"
	"testing"

	"github.com/YUJIAJING0408/go-oidn/oidn"
)

func TestMultiRoundDenoising(t *testing.T) {
	// 初始化库
	oidn.SetLibraryPath("I:\\Codes\\Go\\go-oidn\\lib\\windows")
	oidn.SetVersion("v2.4.1")
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

	// 多轮去噪通常用于相同大小图片进行批量去噪————切片去噪，当图像过大时可以分块进行去噪
	tileColorDatas, tileWidth, tileHeight, w, h, err := oidn.LoadPNGTilesWithOverlap("I:\\Codes\\Go\\go-oidn\\noisy1.png", 4, 4, 16)
	if err != nil {
		panic(err)
	}

	// 创建缓冲区
	colorBuf, _ := device.NewBuffer(tileWidth * tileHeight * 3 * 4)
	defer colorBuf.Release()
	outputBuf, _ := device.NewBuffer(tileWidth * tileHeight * 3 * 4)
	defer outputBuf.Release()

	// 创建滤波器
	filter, err := device.NewFilter(oidn.FilterTypeRT)
	if err != nil {
		panic(err)
	}
	defer filter.Release()

	filter.SetImage(oidn.ImageColor, colorBuf, oidn.FormatFloat3, tileWidth, tileHeight)
	filter.SetImage(oidn.ImageOutput, outputBuf, oidn.FormatFloat3, tileWidth, tileHeight)
	filter.SetBool(oidn.BoolHDR, false)
	filter.SetInt(oidn.IntQuality, int(oidn.QualityHigh))
	filter.Commit()

	tileResults := make([][]float32, len(tileColorDatas))
	fmt.Println("Executing...")
	for i, v := range tileColorDatas {
		fmt.Println("Executing ", i)
		colorBuf.Write(v)
		filter.Execute()
		if err := device.GetError(); err != nil {
			panic(err)
		}
		tileResults[i] = outputBuf.Read(tileWidth * tileHeight * 3)
	}

	// 保存结果
	if err = oidn.SavePNGTilesWithOverlap(tileResults, 4, 4, 16, w, h, "I:\\Codes\\Go\\go-oidn\\denoised1.png"); err != nil {
		return
	}
	fmt.Println("Done.")
}
