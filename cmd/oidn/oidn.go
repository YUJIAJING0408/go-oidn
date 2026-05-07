package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/YUJIAJING0408/go-oidn/oidn"
	"github.com/spf13/cobra"
)

var denoiseCmd = &cobra.Command{
	Use:     "denoise [input.png]",
	Aliases: []string{"dn"},
	Short:   "Denoise an image",
	Args:    cobra.ExactArgs(1),
	RunE:    runDenoise,
}

var (
	outputPath string
	albedoPath string
	normalPath string
	useHDR     bool
	quality    int
)

var rootCmd = &cobra.Command{
	Use:   "oidn",
	Short: "Intel Open Image Denoise CLI",
	Long:  "A command-line tool for denoising images using Intel Open Image Denoise.",
}

func init() {
	// 注册参数
	denoiseCmd.Flags().StringVarP(&outputPath, "output", "o", "denoised.png", "output image path")
	denoiseCmd.Flags().StringVarP(&albedoPath, "albedo", "a", "", "albedo auxiliary image (diffuse color)")
	denoiseCmd.Flags().StringVarP(&normalPath, "normal", "n", "", "normal auxiliary image (world-space, range -1..1)")
	denoiseCmd.Flags().BoolVarP(&useHDR, "hdr", "H", false, "input is HDR (high dynamic range)")
	denoiseCmd.Flags().IntVarP(&quality, "quality", "q", int(oidn.QualityHigh), "quality mode (0=default,4=fast,5=balanced,6=high)")
	rootCmd.AddCommand(denoiseCmd)
}

func runDenoise(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// 初始化 OIDN
	if err := oidn.Init(); err != nil {
		return fmt.Errorf("init OIDN: %w", err)
	}

	// 加载输入图像
	fmt.Printf("Loading %s ...\n", inputPath)
	colorData, width, height, err := oidn.LoadPNG(inputPath)
	if err != nil {
		return fmt.Errorf("load input: %w", err)
	}

	// 加载或生成辅助图像
	var albedoData, normalData []float32

	if albedoPath != "" {
		albedoData, _, _, err = oidn.LoadPNG(albedoPath)
		if err != nil {
			return fmt.Errorf("load albedo: %w", err)
		}
	} else {
		// 生成默认全白 albedo
		albedoData = make([]float32, width*height*3)
		for i := range albedoData {
			albedoData[i] = 0.8
		}
		fmt.Println("  Using default white albedo.")
	}

	if normalPath != "" {
		normalData, _, _, err = oidn.LoadPNG(normalPath)
		if err != nil {
			return fmt.Errorf("load normal: %w", err)
		}
	} else {
		// 生成默认正前法线 (0,0,1)
		normalData = make([]float32, width*height*3)
		for i := 0; i < width*height; i++ {
			normalData[i*3+2] = 1.0
		}
		fmt.Println("  Using default forward normal (0,0,1).")
	}

	// 创建设备（CPU）
	device, err := oidn.NewDevice(oidn.DeviceTypeCPU)
	if err != nil {
		return fmt.Errorf("create device: %w", err)
	}
	defer device.Release()

	if err := device.Commit(); err != nil {
		return fmt.Errorf("commit device: %w", err)
	}

	// 创建缓冲区
	colorBuf, _ := device.NewBuffer(len(colorData) * 4)
	defer colorBuf.Release()
	colorBuf.Write(colorData)

	albedoBuf, _ := device.NewBuffer(len(albedoData) * 4)
	defer albedoBuf.Release()
	albedoBuf.Write(albedoData)

	normalBuf, _ := device.NewBuffer(len(normalData) * 4)
	defer normalBuf.Release()
	normalBuf.Write(normalData)

	outputBuf, _ := device.NewBuffer(width * height * 3 * 4)
	defer outputBuf.Release()

	// 创建 RT 滤波器
	filter, err := device.NewFilter(oidn.FilterTypeRT)
	if err != nil {
		return fmt.Errorf("create filter: %w", err)
	}
	defer filter.Release()

	// 设置滤波器参数
	filter.SetImage(oidn.ImageColor, colorBuf, oidn.FormatFloat3, width, height)
	filter.SetImage(oidn.ImageAlbedo, albedoBuf, oidn.FormatFloat3, width, height)
	filter.SetImage(oidn.ImageNormal, normalBuf, oidn.FormatFloat3, width, height)
	filter.SetImage(oidn.ImageOutput, outputBuf, oidn.FormatFloat3, width, height)
	filter.SetBool(oidn.BoolHDR, useHDR)
	filter.SetInt(oidn.IntQuality, quality)

	filter.Commit()

	// 执行降噪
	fmt.Println("Denoising ...")
	filter.Execute()
	if err := device.GetError(); err != nil {
		return fmt.Errorf("denoise failed: %w", err)
	}

	// 保存结果
	result := outputBuf.Read(width * height * 3)
	absOutput, _ := filepath.Abs(outputPath)
	if err := oidn.SavePNG(result, width, height, outputPath); err != nil {
		return fmt.Errorf("save output: %w", err)
	}
	fmt.Printf("Saved to %s\n", absOutput)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
