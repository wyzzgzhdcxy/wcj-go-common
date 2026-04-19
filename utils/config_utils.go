package utils

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

// LoadPropertiesConfig loadConfig 加载指定路径的 .properties 文件并返回 *properties.Properties 对象
func LoadPropertiesConfig(filePath string) (*properties.Properties, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 使用 properties 加载文件
	props, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %v", err)
	}

	return props, nil
}
