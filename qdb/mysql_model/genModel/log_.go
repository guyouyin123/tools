package genModel

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	stdLog "log"
	"os"
	"time"
)

// StdLog 本地终端日志输出
func StdLog() logger.Interface {
	newLogger := logger.New(
		stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	return newLogger
}

// LogPath 指定日志路径输出
func LogPath(logPath string) logger.Interface {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("无法打开日志文件")
	}
	newLogger := logger.New(
		stdLog.New(logFile, "\r\n", stdLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	return newLogger
}

// LogOther 使用外接其他日志
func LogOther() logger.Interface {
	newLogger := CustomLogger{logger.Default.LogMode(logger.Info)}
	return newLogger
}

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
	sql, rows := fc() // 获取 SQL 查询和受影响的行数
	duration := time.Since(begin)
	logInfo := fmt.Sprintf("SQL: %s | Duration: %s | Rows affected: %d", sql, duration, rows)
	fmt.Println(logInfo) //TODO 替换为其他日志输出
}
