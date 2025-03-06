/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package main

import (
	"example/config"
	"example/entity"
	"example/route"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ReadConfig()
	// 初始化mysql
	config.InitMysqlDataBase()
	// 初始化redis
	config.InitRedis()
	// 自动迁移数据库
	config.MysqlDataBase.AutoMigrate(&entity.User{})
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.Gin.CorsAllowOrigins, // 允许的前端来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour, // 预检请求的缓存时间
	}))
	// app.Static("/api/uploads", "./uploads") 静态文件访问路径以及目录
	route.RegisterRoutes(app)
	app.Run(config.Config.Gin.Port)
}
