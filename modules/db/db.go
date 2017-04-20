package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"script-web/modules/config"
	"fmt"
)

var (
	//Host     string = "127.0.0.1:3306"
	//Database string = "script-web"
	//User     string = "root"
	//Password string = "root"
	db       *gorm.DB
	mysqldb  *MysqlDB
)


type MysqlDB struct {
	Db *gorm.DB
	Name string
}

func GetDb() (*MysqlDB, error) {
	if db == nil {
		var err error
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", config.MysqlUser, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase))
		// "root:root@tcp/cloudcli?charset=utf8&parseTime=True&loc=Local"
		if err != nil {
			return nil, err
		}
	}

	if mysqldb == nil {
		mysqldb = &MysqlDB{
			Db: db,
			Name: "mysqldb",
		}
	}

	return mysqldb, nil
}

func (this *MysqlDB) Close() error {
	return this.Db.Close()
}