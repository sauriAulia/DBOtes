package models

import "gorm.io/gorm"

func MigrateModel(db *gorm.DB) error {
	return db.AutoMigrate(&Customer{}, &Order{}, &User{})
}
