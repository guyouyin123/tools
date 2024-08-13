package qdateTime

import "time"

// GetZeroTime 获取时间的零点时间
func GetZeroTime(timestamp int64) time.Time {
	t := time.Unix(timestamp, 0)
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return zeroTime
}

// GetDates 时间分布--取时间范围的点(以天为单位)，首位一个点，中间点平均分布,如不均匀则更靠近endTime的时间
// number 采点的数量
func GetDates(startTime, endTime int64, number int) []string {
	// 计算时间差,单位是秒
	timeDiff := endTime - startTime

	// 如果时间差小于等于6天,则返回全部时间范围内的点
	if timeDiff <= (int64(number)-1)*86400 {
		dates := make([]string, 0, timeDiff/86400+1)
		for t := startTime; t <= endTime; t += 86400 {
			dates = append(dates, time.Unix(t, 0).Format("2006-01-02"))
		}
		return dates
	}

	// 生成number个点的时间戳
	timestamps := make([]int64, number)
	timestamps[0] = startTime
	timestamps[number-1] = endTime
	step := (timeDiff / 86400) / (int64(number) - 1)
	for i := number - 2; i > 0; i-- {
		timestamps[i] = endTime - (step*86400)*int64(number-1-i)
	}
	// 转换为日期字符串
	dates := make([]string, number)
	for i, t := range timestamps {
		dates[i] = time.Unix(t, 0).Format("2006-01-02")
	}
	return dates
}

// GetRecentWeekDates 获取当前时间前n个礼拜时间段--周天-周六
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
