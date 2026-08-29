package core

import (
	"fmt"
	"strconv"
	"time"

	"log"

	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
)

func GetCurTime() string {
	curTime := time.Now() //2019-07-31 13:55:21.3410012 +0800 CST m=+0.006015601
	return curTime.Format("20060102150405")
}

func GetTime() string {
	curTime := time.Now() //2019-07-31 13:55:21.3410012 +0800 CST m=+0.006015601
	return curTime.Format("15:04:05")
}

func GetCurDate() string {
	curTime := time.Now() //2019-07-31 13:55:21.3410012 +0800 CST m=+0.006015601
	return curTime.Format("20060102")
}

func LastNDay(n int) time.Time {
	now := time.Now()
	// 获取前一天的时间
	return now.AddDate(0, 0, -n)
}

func Time2Str(time time.Time) string {
	return time.Format("20060102")
}

func Time2DateStr(time time.Time) string {
	return time.Format("20060102")
}

func GetYesterdayDateStr() string {
	// 获取当前时间
	now := time.Now()
	// 获取昨天的日期（通过减去24小时可能不是最佳方法，因为有时区变化和夏令时的影响）
	// 昨天 = 现在 - 1天
	yesterday := now.AddDate(0, 0, -1)

	// 格式化日期为字符串（例如：2006-01-02）
	return yesterday.Format("20060102")
}

func UpdateSystemDate(dateTime string) bool {
	_, err1 := gproc.ShellExec(`date  ` + gstr.Split(dateTime, " ")[0])
	_, err2 := gproc.ShellExec(`time  ` + gstr.Split(dateTime, " ")[1])
	if err1 != nil && err2 != nil {
		log.Printf("更新系统时间错误:请用管理员身份启动程序!")
		return false
	} else {
		log.Printf("更新成功")
		return true
	}
}

func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func DateFormat(time time.Time) string {
	return time.Format("20060102")
}

func StringToTime(timeStr string) time.Time {
	curTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Printf("结束时间字符串异常:%s", timeStr)
	}
	return curTime
}

func DateStringToTime(dateStr string) time.Time {
	curTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		log.Printf("结束时间字符串异常:%s", dateStr)
	}
	return curTime
}

// weekdayChineseMap 定义了星期几到中文的映射
var weekdayChineseMap = map[time.Weekday]string{
	time.Sunday:    "日",
	time.Monday:    "一",
	time.Tuesday:   "二",
	time.Wednesday: "三",
	time.Thursday:  "四",
	time.Friday:    "五",
	time.Saturday:  "六",
}

// WeekdayToChinese  将星期几转换为中文
func WeekdayToChinese(weekday time.Weekday) string {
	return weekdayChineseMap[weekday]
}

func TimeDifference(timeStr string) int {
	// 定义两个时间点
	layout := "20060102"
	// 解析时间字符串为 time.Time 类型
	time1, err1 := time.Parse(layout, timeStr)
	if err1 != nil {
		fmt.Println("Error parsing time1:", err1)
		return 0
	}
	// 计算时间差
	duration := time.Now().Sub(time1)
	// 将时间差转换为天数
	// 注意：这里假设一天总是24小时，不考虑闰秒等复杂情况
	return int(duration.Hours() / 24)
}

func TimeDifferenceWithTime(lastTime time.Time) int {
	// 计算时间差
	duration := time.Now().Sub(lastTime)
	// 将时间差转换为天数
	// 注意：这里假设一天总是24小时，不考虑闰秒等复杂情况
	return int(duration.Hours() / 24)
}

func CalculateDaysDifference(dateStr string) (error, int) {
	parsedDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return err, 0
	}
	// 标准化为当天的 00:00:00（忽略时分秒）
	parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	// 计算时间差
	duration := time.Now().Sub(parsedDate)
	// 将时间差转换为天数
	// 注意：这里假设一天总是24小时，不考虑闰秒等复杂情况
	return nil, int(duration.Hours() / 24)
}

func ToYesterdayStr(dateStr string) string {
	// 解析日期
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		fmt.Println("解析日期出错:", err)
		return ""
	}
	// 减去一天
	previousDay := date.AddDate(0, 0, -1)
	// 格式化为相同的字符串格式
	result := previousDay.Format("20060102")
	return result
}

func GenerateTimeRange(startStr, endStr, layout string, step time.Duration) ([]string, error) {
	// 解析时间字符串
	start, err := time.ParseInLocation(layout, startStr, time.Local)
	if err != nil {
		return nil, fmt.Errorf("解析起始时间失败: %v", err)
	}
	end, err := time.ParseInLocation(layout, endStr, time.Local)
	if err != nil {
		return nil, fmt.Errorf("解析结束时间失败: %v", err)
	}

	// 检查时间范围
	if end.Before(start) {
		return nil, fmt.Errorf("结束时间必须晚于起始时间")
	}
	// 生成中间时间点
	var result []string
	for current := start; !current.After(end); current = current.Add(step) {
		result = append(result, current.Format(layout))
	}
	return result, nil
}

// GenerateDaysRange start - 20250101     end -20250103
// 输出 [20250101 20250102 20250103]
func GenerateDaysRange(start, end string) ([]string, error) {
	layout := "20060102"
	step := time.Hour * 24 // 按天递增
	return GenerateTimeRange(start, end, layout, step)
}

// 给日期字符串加一天
func AddOneDay(dateStr string) (string, error) {
	// 1. 定义日期格式：Go 使用特定参考时间 "20060102" 表示 YYYYMMDD
	const layout = "20060102"
	// 2. 将字符串解析为 time.Time
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", fmt.Errorf("解析日期失败: %v", err)
	}
	// 3. 加一天（推荐使用 AddDate，更可靠，尤其是跨月跨年时）
	t = t.AddDate(0, 0, 1) // 年, 月, 日
	// 4. 格式化回原来的字符串格式
	newDateStr := t.Format(layout)
	return newDateStr, nil
}
func Timestamp2Time(timestampStr string) (time.Time, error) {
	// 使用 strconv.ParseInt 转为 int64
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		fmt.Println("字符串转 int64 失败:", err)
		return time.Now(), err
	}
	// 转为秒和纳秒
	sec := timestamp / 1000          // 秒
	nsec := (timestamp % 1000) * 1e6 // 毫秒转纳秒
	// 构造 time.Time
	return time.Unix(sec, nsec), nil
}
