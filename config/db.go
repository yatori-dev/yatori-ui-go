package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"yatory-gui-wails3/model"
)

var DB *gorm.DB

func InitDB() {
	_, err2 := os.Stat("C:\\db")
	if err2 != nil {
		os.Mkdir("C:\\db", os.ModePerm)
	}
	_db, err := gorm.Open(sqlite.Open("C:\\db\\yatori.db"))
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	err = _db.AutoMigrate(&model.XueXiTong{})
	if err != nil {
		panic("failed to migrate database")
	}
	DB = _db
}
