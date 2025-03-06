/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDataBase *gorm.DB

func InitMysqlDataBase() {
	dbc := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.MySQL.User, Config.MySQL.Password, Config.MySQL.Host, Config.MySQL.Port, Config.MySQL.Database)
	db, err := gorm.Open(mysql.Open(dbc), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库时发生错误:%v", err)
		os.Exit(1)
	}
	MysqlDataBase = db
	fmt.Println("连接数据库成功")
}
