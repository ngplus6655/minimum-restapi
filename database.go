package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Database struct{
	Service string
	User    string
	Pass    string
	DatabaseName string
}

func (d Database) connect() (*gorm.DB, error) {
	connStr := d.User + ":" + d.Pass + "@/" + d.DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(d.Service, connStr)
	return db, err
}

func (d Database) init() *gorm.DB {
	db, err := d.connect()
	if err != nil {
		log.Fatalln("データベースの接続に失敗しました。")
	}
	return db
}

func (d Database) migrate() *gorm.DB {
	db, err := d.connect()
	if err != nil {
		log.Fatalln("データベースの接続に失敗しました。")
	}
	db.AutoMigrate(&Article{})
	return db
}