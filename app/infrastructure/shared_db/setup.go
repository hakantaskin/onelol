package shared_db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Name string
	Host string
	Port string
	Pass string
	DB   string
}

func NewConnector(config DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Name, config.Pass, config.Host, config.Port, config.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ProfileEntry{},
	)
}
