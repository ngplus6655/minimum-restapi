package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// testを実行する関数型
type test func(db *gorm.DB)

func transactionTestArticles(db *gorm.DB, t test) error {
	tx := db.Begin()
	// test実行
	t(tx)	
	err := tx.Rollback().Error
	return err
}

