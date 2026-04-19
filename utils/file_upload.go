package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"wcj-go-common/core"
)

func UploadFile2TencentCloud(filename string) error {
	targetURL := "http://111.229.201.94:10030/upload"
	// 打开要上传的文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer core.Close(file)

	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建表单文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("创建表单字段失败: %v", err)
	}

	// 将文件内容复制到表单
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 关闭writer以完成表单构建
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("关闭writer失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置Content-Type头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer core.Close(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上传失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("上传成功! 服务器响应: %s\n", string(respBody))
	return nil
}
