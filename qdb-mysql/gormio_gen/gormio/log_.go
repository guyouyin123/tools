package gormio

import (
	"context"
	"gorm.io/gorm/logger"
	"time"
)

// CustomLogger 用于捕获 SQL 查询
type CustomLogger struct {
	logger.Interface
}

// LogMode 设置日志级别
func (c CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return CustomLogger{c.LogMode(level)}
}

// Trace 用于捕获 SQL 查询
func (c CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	//sql, rows := fc() // 获取 SQL 查询和受影响的行数
	//duration := time.Since(begin)
	//log.Info("gorm", "data", fmt.Sprintf("SQL: %s | Duration: %s | Rows affected: %d", sql, duration, rows))
}
