package main

import (
	"os"
	"github.com/jinzhu/gorm"
)

func connTestDB() *gorm.DB {
	testdb := Database{
		Service: "mysql",
		User: os.Getenv("MINIMUM_APP_DATABASE_USER"),
		Pass: os.Getenv("MINIMUM_APP_DATABASE_PASS"),
		DatabaseName: os.Getenv("MINIMUM_APP_TEST_DATABASE_NAME"),
	}
	db := testdb.migrate()
	return db
}

func setFixture() *gorm.DB {
	db := connTestDB()
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


func cleanUpFixture(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE articles;")
	db.Close()
}


func fetchTestDB() Database {
	testdb := Database{
		Service: "mysql",
		User: os.Getenv("MINIMUM_APP_DATABASE_USER"),
		Pass: os.Getenv("MINIMUM_APP_DATABASE_PASS"),
		DatabaseName: os.Getenv("MINIMUM_APP_TEST_DATABASE_NAME"),
	}
	return testdb
}