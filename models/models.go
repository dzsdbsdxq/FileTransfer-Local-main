package models

import (
	"FileTransfer/pkg/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Model struct {
	ID    int `gorm:"primary_key" json:"id"`
	Ctime int `json:"ctime"`
}

func init() {
	var err error

	db, err = gorm.Open("sqlite3", "./data/data.db")

	if err != nil {
		logging.Logger.Errorf("%v\n", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tb_" + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

}
