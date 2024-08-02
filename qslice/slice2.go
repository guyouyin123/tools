package qslice

import "time"

// 获取当前时间前n个礼拜时间段--周天-周六
func GetRecentWeekDates(n int) [][]string {
	// 获取当前时间
	now := time.Now()
	// 初始化结果切片
	result := make([][]string, 0, n)
	// 循环最近n个礼拜
	for i := 0; i < n; i++ {
		// 获取当前礼拜的第一天和最后一天
		firstDay := now.AddDate(0, 0, -(int(now.Weekday())+7)%7-n*7)
		lastDay := firstDay.AddDate(0, 0, 6)
		// 将当前礼拜的起止日期添加到结果切片
		result = append(result, []string{firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")})
		// 移到下一个礼拜
		now = now.AddDate(0, 0, 7)
	}
	return result
}
