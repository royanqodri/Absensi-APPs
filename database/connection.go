package database

import (
	"fmt"

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.Get().DBHost,
		config.Get().DBUsername,
		config.Get().DBPassword,
		config.Get().DBName,
		config.Get().DBPort,
	)
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	sqlDB, err := DBConn.DB()
	util.PanicIfError(err)

	// Ping the database to confirm a successful connection
	if err = sqlDB.Ping(); err != nil {
		util.PanicIfError(err)
	}

	// Force apply statement_timeout to every new connection in pool
	sqlDB.Exec("SET statement_timeout = '15s'") // <- will apply only to this session
}
