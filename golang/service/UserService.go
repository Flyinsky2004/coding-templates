package service

import (
	"example/config"
	"example/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyInfo(c *gin.Context) {
	var user entity.User
	userId, _ := c.Get("userId")
	if err := config.MysqlDataBase.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "在获取用户信息时出错！详细信息:"+err.Error()))
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, entity.SuccessResponse(user))
}
