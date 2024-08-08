package genModel

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	stdLog "log"
	"os"
	"time"
)

var (
	WriteDB *gorm.DB
	// 读库，预留 ReadDB  *gorm.DB
)

func GetDbUrl() (string, error) {
	//配置
	dataSource := map[string]interface{}{}
	ip := dataSource["ip"].(string)
	port := dataSource["port"].(int)
	username := dataSource["username"].(string)
	password := dataSource["password"].(string)
	dbname := dataSource["dbname"].(string)
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC", username, password, ip, port, dbname)
	return dbSource, nil
}

func InitDb() error {
	dbUrl, err := GetDbUrl()
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
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbUrl, // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})

	if err != nil {
		panic(fmt.Sprintf("db gorm err:%v", err))
	}
	WriteDB = db
	return nil
}

// NewCtx 新生成连接
func NewCtx(ctx context.Context) *gorm.DB {
	return WriteDB.WithContext(ctx)
}

// Transaction 事务执行
func Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	db := NewCtx(ctx)
	return db.Transaction(fn)
}
