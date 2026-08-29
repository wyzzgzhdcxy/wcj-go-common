package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/wyzzgzhdcxy/wcj-go-common/core"
)

// InitLog 初始化日志 all -true 包含文件和控制台  false-仅仅控制台
func InitLog(all bool) {
	// 1. 打开日志文件（如果不存在则创建，追加写入）
	logFn := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	file, err := os.OpenFile(core.GetTempDir()+"/"+logFn+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err) // 如果连日志文件都打不开，程序无法运行
	}
	// 2. 创建一个 MultiWriter：同时输出到 os.Stdout（控制台）和 file（日志文件）
	// 注意：打包后的 GUI 程序没有控制台，os.Stdout.Fd() 返回 0，此时只用文件
	if all && os.Stdout.Fd() != 0 {
		multiWriter := io.MultiWriter(os.Stdout, file)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(file)
	}
}
