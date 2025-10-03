package data

import (
	"context"
	"time"

	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data struct definition
type Data struct {
	db *gorm.DB
}

func (d *Data) WithContext(ctx context.Context) {
	panic("unimplemented")
}

func NewData(s *conf.Data) (*Data, func(), error) {
	// Initialize your Data struct here using the MySQL configuration
	db, err := gorm.Open(mysql.Open(s.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// 判断数据库连接是否成功
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	// 设置数据库链接池参数
	sqlDB.SetMaxIdleConns(s.MySQL.MaxidleConns)                                // 设置空闲连接池的最大连接数
	sqlDB.SetMaxOpenConns(s.MySQL.MaxOpenConns)                                // 设置数据库的最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(s.MySQL.MaxLifetime) * time.Second) // 设置连接的最大可复用时间

	// Return a cleanup function to close the database connection
	sqlcleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return &Data{
		db: db,
	}, sqlcleanup, nil
}
