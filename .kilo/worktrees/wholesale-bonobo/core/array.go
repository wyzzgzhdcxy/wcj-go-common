package core

import (
	"strconv"
	"strings"
)

// ArrayContains contains 判断切片中是否包含某个元素
func ArrayContains[T comparable](slice []T, elem T) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// DimensionalInterface2String 2维数组转字符串
func DimensionalInterface2String(list [][]interface{}) [][]string {
	var resultArr [][]string
	for _, item1 := range list {
		var arr []string
		for _, item2 := range item1 {
			str := ToString(item2)
			arr = append(arr, Base64DecodeStr(str[1:len(str)-1]))
		}
		resultArr = append(resultArr, arr)
	}
	return resultArr
}

func DimensionalToStringResult(dimensional *[][]string) StringResult {
	strResult := StringResult{}
	strResult.Count = 0
	var builder strings.Builder
	for _, arr := range *dimensional {
		for _, item := range arr {
			strResult.Count++
			builder.WriteString(item + "\r\n")
		}
	}
	strResult.Result = builder.String()
	return strResult
}

type StringResult struct {
	Result string
	Count  int
}

// IntArrayToString []int{1, 2, 3, 4, 5}
// 输出: "1, 2, 3, 4, 5"
func IntArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = strconv.Itoa(v)
	}
	return strings.Join(strArr, ", ")
}

// CountStringOccurrences 统计一个数组中各字符串出现的次数
func CountStringOccurrences(arr []string) map[string]int {
	// 创建一个空的 map 来存储字符串及其出现次数
	counts := make(map[string]int)

	// 遍历数组中的每个字符串
	for _, str := range arr {
		// 将字符串出现的次数加 1
		counts[str]++
	}

	return counts
}

// SplitStringSlice 将一个字符串切片切分成指定长度的多个子切片
func SplitStringSlice(slice []string, chunkSize int) [][]string {
	var result [][]string
	length := len(slice)
	for i := 0; i < length; i += chunkSize {
		end := i + chunkSize
		if end > length {
			end = length
		}
		result = append(result, slice[i:end])
	}
	return result
}

// Difference returns the elements in a that are not in b
func Difference(a, b *[]string) *[]string {
	// 将数组 b 转换为 map，以便快速查找
	bMap := make(map[string]struct{})
	for _, val := range *b {
		bMap[val] = struct{}{}
	}

	var diff []string
	for _, val := range *a {
		if _, exists := bMap[val]; !exists {
			diff = append(diff, val)
		}
	}
	return &diff
}

// RemoveDuplicateStrings 接收一个字符串切片，并返回一个新的、已去重的切片。
func RemoveDuplicateStrings(strings []string) []string {
	// 创建一个map来追踪已经遇到的字符串。
	seen := make(map[string]bool)
	// 创建一个新的切片来存储去重后的结果。
	var result []string

	// 遍历原始切片。
	for _, str := range strings {
		// 如果当前字符串尚未在map中出现过，则将其添加到结果切片中，并在map中标记为已见。
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	// 返回去重后的切片。
	return result
}
