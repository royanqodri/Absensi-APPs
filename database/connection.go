package database

import (
	"fmt"
	"time"

	"github.com/royanqodri/Absensi/config"
	"github.com/royanqodri/Absensi/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConn          *gorm.DB
	DBConnOperation *gorm.DB
)

func ConnectDatabase() {
	var err error

	// ✅ PERBAIKI: Gunakan format DSN yang benar untuk MySQL
	// Format: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get().DBUsername,
		config.Get().DBPassword,
		config.Get().DBHost,
		config.Get().DBPort,
		config.Get().DBName,
	)

	// Log DSN (tanpa password untuk keamanan)
	fmt.Printf("🔗 Connecting to MySQL: %s@%s:%d/%s\n",
		config.Get().DBUsername,
		config.Get().DBHost,
		config.Get().DBPort,
		config.Get().DBName)

	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	sqlDB, err := DBConn.DB()
	util.PanicIfError(err)

	// Ping the database to confirm a successful connection
	if err = sqlDB.Ping(); err != nil {
		util.PanicIfError(err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Database connected successfully")
}
