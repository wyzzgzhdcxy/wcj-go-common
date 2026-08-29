package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ReadCSVFile(filepath string) ([][]string, int) {
	// 打开 CSV 文件
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("无法打开文件: %v", err)
		panic(err)
	}
	defer Close(file)
	reader := csv.NewReader(file)
	var records [][]string
	var errLineCount int
	// 创建一个新的 CSV 读取器
	// 读取 CSV 文件中的所有记录
	// 逐行读取 CSV 文件
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // 文件读取完毕
		}
		if err != nil {
			// 忽略错误行并继续读取下一行
			errLineCount++
			continue
		}
		records = append(records, record)
	}
	return records, errLineCount
}

// ReadCSVFile2 读取csv数据，如果数据有错误,直接异常
func ReadCSVFile2(filepath string, gbk bool) [][]string {
	// 打开 CSV 文件
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("无法打开文件: %v", err)
	}
	defer Close(file)

	var reader *csv.Reader
	// 创建一个 GBK 到 UTF-8 的转换器
	if gbk {
		gbkReader := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())
		reader = csv.NewReader(gbkReader)
	} else {
		reader = csv.NewReader(file)
	}

	// 创建一个新的 CSV 读取器
	// 读取 CSV 文件中的所有记录
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("读取CSV文件失败: %v", err)
	}
	return records
}
