package core

import (
	"fmt"
)

func PrintProgressBarDetail(current, total int) {
	width := 50 // 进度条宽度
	percent := float64(current) / float64(total) * 100
	filledWidth := int(float64(width) * percent / 100)

	// 构建进度条
	bar := "["
	for i := 0; i < width; i++ {
		if i < filledWidth {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"
	// 打印进度条和百分比
	fmt.Printf("\r%s %.2f%%", bar, percent)
}

// PrintProgressBar 简化版,只打印整数倍进度,打印的过于详细，经过测试影响性能
func PrintProgressBar(current, total int) {
	step := total / 100
	if total%100 != 0 {
		step = step + 1
	}
	if current%step != 0 {
		//如果不是整数倍无需打印
		return
	}
	PrintProgressBarDetail(current, total)
}

func PrintProgressBarFinish() {
	PrintProgressBar(100, 100)
	fmt.Println()
}
