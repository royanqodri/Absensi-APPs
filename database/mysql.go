package database

import (
	"github.com/royanqodri/Absensi/config"

	"github.com/royanqodri/Absensi/model/entity"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql(cfg *config.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	DB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return DB
}

// InitialMigration - Auto migrate tabel
func InitialMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.TAbsensi{},
		&entity.TUser{},
	)
	if err != nil {
		panic("Failed to migrate: " + err.Error())
	}
	fmt.Println("✅ Migration completed successfully")
}

func MigrateDatabase() {
	if DBConn == nil {
		panic("❌ Database not connected. Call ConnectDatabase() first")
	}
	InitialMigration(DBConn)
}
