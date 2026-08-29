package core

import (
	"fmt"
	"strings"
)

// MySet Set 表示一个集合
type MySet struct {
	data map[interface{}]struct{}
}

// NewSet 创建一个新的 Set
func NewSet() *MySet {
	return &MySet{
		data: make(map[interface{}]struct{}),
	}
}

// Add 添加一个元素到 Set
func (s *MySet) Add(item interface{}) {
	s.data[item] = struct{}{}
}

// Remove 从 Set 中移除一个元素
func (s *MySet) Remove(item interface{}) {
	delete(s.data, item)
}

// Contains 检查 Set 是否包含某个元素
func (s *MySet) Contains(item interface{}) bool {
	_, exists := s.data[item]
	return exists
}

// Size 返回 Set 的大小（元素数量）
func (s *MySet) Size() int {
	return len(s.data)
}

// String 将 Set 转换为字符串表示形式
func (s *MySet) String() string {
	items := s.ToArray()
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}

func (s *MySet) ToArray() []string {
	items := make([]string, 0, len(s.data))
	for item := range s.data {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return items
}
