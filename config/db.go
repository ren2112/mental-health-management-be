package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mental-health-management-be/models"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	// 修改为你自己的数据库信息
	username := "root"
	password := "123456"
	host := "127.0.0.1"
	port := "3306"
	dbname := "mental_health_management"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	// ✅ 自定义 GORM 日志配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,
			Colorful:                  true, // 彩色日志
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	sqlDB, _ := db.DB()

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	log.Println("数据库连接成功")
	err = DB.AutoMigrate(
		&models.Student{},
		&models.Teacher{},
		&models.Manager{},
		&models.Post{},
		&models.Appointment{},
	)

	if err != nil {
		log.Fatal("数据表创建失败:", err)
	}

	log.Println("数据表迁移完成")
}
