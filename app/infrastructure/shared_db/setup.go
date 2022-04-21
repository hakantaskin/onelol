package shared_db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnector() *gorm.DB {
	dsn := "onelol:hkn1242@tcp(139.59.147.81:3306)/onelol?charset=utf8&parseTime=True&loc=Local"
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
