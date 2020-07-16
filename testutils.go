package main

import (
	"os"
	"log"

	"github.com/jinzhu/gorm"
)

func SetFixture() *gorm.DB {
	dbname = os.Getenv("MINIMUM_APP_TEST_DATABASE_NAME")
	d := Database{
		Service:      dbservice,
		User:         dbuser,
		Pass:         dbpass,
		DatabaseName: dbname,
	}
	db, err := d.connect()
	if err != nil {
		log.Fatalln("データベースの接続に失敗しました。")
	}
	db.AutoMigrate(&Article{})

	articles := Articles{
		Article{Title: "test1", Desc: "test description1", Content: "test content1"},
		Article{Title: "test2", Desc: "test description2", Content: "test content2"},
		Article{Title: "test1", Desc: "test description3", Content: "test content3"},
	}
	for _, article := range articles {
		db.Create(&article)
	}
	return db
}



func CleanUpFixture(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE articles;")
	db.Close()
}