package core

import "regexp"

// MatchFirstGroupString 通过正则表达式匹配字符串,获取到的第一条匹配记录的第一个分组数据
func MatchFirstGroupString(input string, regexpString string) string {
	regExpCompile := regexp.MustCompile(regexpString)
	matches := regExpCompile.FindAllStringSubmatch(input, -1)
	// 打印所有匹配的电话号码
	for _, match := range matches {
		// match[0] 是整个匹配的字符串（包括捕获组）
		// match[1] 是第一个捕获组（即我们想要的电话号码）
		return match[1]
	}
	return ""
}
