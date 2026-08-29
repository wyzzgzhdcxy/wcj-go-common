package core

import (
	"encoding/json"
	"fmt"
	"github.com/go-basic/uuid"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func StrIsEmpty(str string) bool {
	return len(str) == 0
}
func StrIsNotEmpty(str string) bool {
	return len(str) != 0
}

// ContainsChinese 判断字符串是否包含中文
func ContainsChinese(str string) bool {
	for _, runeValue := range str {
		if unicode.Is(unicode.Han, runeValue) {
			return true
		}
	}
	return false
}

// ContainsJapanese 包含日文字符
// isJapanese checks if the input string contains any Japanese characters.
func ContainsJapanese(s string) bool {
	for _, r := range s {
		// Unicode blocks for Japanese characters:
		// Hiragana, Katakana, Kanji, Katakana Phonetic Extensions, etc.
		switch {
		case unicode.Is(unicode.Hiragana, r):
			return true
		case unicode.Is(unicode.Katakana, r):
			return true
			// You can add more specific Japanese-related Unicode blocks if needed
		}
	}
	return false
}

func PadStringWithZeros(input string, length int) string {
	// 使用 fmt.Sprintf 格式化字符串，%0<length>s 表示用0填充到指定长度，s表示字符串
	// 注意：这里 input 应该是可以转换为字符串的任何类型，但在这个例子中我们假设它是字符串
	// 如果 input 是一个整数，您可以使用 strconv.Itoa 将其转换为字符串
	paddedString := fmt.Sprintf("%0*s", length, input)
	// 但是上面的代码有一个问题，它会把 input 当作普通字符串来处理，不会考虑其作为数字的长度
	// 如果我们想要根据数字的长度来填充，我们需要先转换 input 为数字，然后再转回字符串（如果 input 已经是字符串表示的数字）
	// 或者，我们可以直接计算 input 的长度，并用零填充到目标长度
	// 下面的代码演示了如何正确地根据数字字符串的长度来填充零

	// 更正后的代码：如果 input 是数字字符串，我们计算其长度并填充
	inputLen := len(input)
	if inputLen >= length {
		// 如果 input 已经足够长，直接返回（这里可以根据需求决定是否截断）
		return input[:length] // 如果需要截断到指定长度
	} else {
		// 计算需要填充的零的数量
		zeros := length - inputLen
		// 使用 strings.Repeat 来生成指定数量的零字符串
		paddedString = strings.Repeat("0", zeros) + input
	}
	return paddedString
}

// Contains 检查切片是否包含某个字符串
func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func TraceId() string {
	traceId := GetCurTime() + uuid.New()[0:8]
	return traceId
}

// ToString interface 转 string
func ToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case time.Time:
		t, _ := value.(time.Time)
		key = t.String()
		key = strings.Replace(key, " +0800 CST", "", 1)
		key = strings.Replace(key, " +0000 UTC", "", 1)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("转换失败：%s", err)
		panic("转换失败")
	}
	return num
}

// TruncateString 截取字符串指定字数，如果长度不够返回所有字符串
func TruncateString(str string, num int) string {
	if len(str) <= num {
		return str
	}
	return str[:num] + "......"
}

// JaccardSimilarity 判断两个字符串相似度
func JaccardSimilarity(a, b string) float64 {
	setA := make(map[rune]struct{})
	setB := make(map[rune]struct{})

	for _, char := range a {
		setA[char] = struct{}{}
	}
	for _, char := range b {
		setB[char] = struct{}{}
	}

	intersection := 0
	for char := range setA {
		if _, found := setB[char]; found {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 1.0 // 避免除以零
	}
	return float64(intersection) / float64(union)
}

// AssertStrLen 根据字符串长度验证字符串是否符合规定 assertLenStr: "=,100" input长度如果等于100返回true
func AssertStrLen(input string, assertLenStr string) bool {
	if len(assertLenStr) == 0 {
		return true
	}
	i := strings.Index(assertLenStr, ",")
	op := (assertLenStr)[0:i]
	length := StrToInt((assertLenStr)[i+1:])
	if (op == "=" && len(input) == length) || (op == "!=" && len(input) != length) || (op == ">" && len(input) > length) ||
		(op == "<" && len(input) < length) {
		return true
	}
	return false
}
