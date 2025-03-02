package service

import (
	"context"
	"example/config"
	"example/entity"
	"example/utility"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var body UserLoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}
	var user entity.User
	if err := config.MysqlDataBase.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "用户名不存在"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "用户名或密码错误"))
		return
	}
	var token, err = utility.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "生成用户token时发生错误，请再尝试一次，或联系管理员"))
		return
	}
	c.JSON(200, entity.SuccessResponse(token))
}

// SaveCodeToRedis 保存验证码到 Redis
func SaveCodeToRedis(key, code string, ttl time.Duration) error {
	return config.RedisClient.Set(context.Background(), key, code, ttl).Err()
}

// GetCodeFromRedis 获取验证码
func GetCodeFromRedis(key string) (string, error) {
	result, err := config.RedisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("验证码不存在或已过期")
	}
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %v", err)
	}
	return result, nil
}
func DeleteCodeToRedis(key string) error {
	return config.RedisClient.Del(context.Background(), key).Err()
}

// CheckIfCodeExists 检查验证码是否存在
func CheckIfCodeExists(key string) (bool, error) {
	cmd := config.RedisClient.Exists(context.Background(), key)
	exists, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

type RegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}
type SendVerifyCodeRequestBody struct {
	Email string `json:"email"`
}

func SendVerifyCode(c *gin.Context) {
	var reqBody SendVerifyCodeRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "请求表单格式有误"+err.Error()))
		return
	}
	isHasExisted, err := CheckIfCodeExists(reqBody.Email)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if isHasExisted {
		c.JSON(200, entity.ErrorResponse[string](400, "你已经发送了一条未失效的验证码"))
		return
	}
	code := utility.GenerateCode(6)
	if err := SaveCodeToRedis(reqBody.Email, code, time.Minute*3); err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "缓存验证码时发生错误"))
		return
	}

	// 邮件HTML模板
	emailTemplateFront := `
<!DOCTYPE html>
<html dir="ltr" lang="en">
<head>
  <meta name="viewport" content="width=device-width"/>
  <meta content="text/html; charset=UTF-8" http-equiv="Content-Type"/>
  <meta name="color-scheme" content="light"/>
  <meta name="supported-color-schemes" content="light"/>
  <style>
    @font-face {
      font-family: 'Inter';
      font-style: normal;
      font-weight: 400;
      mso-font-alt: 'sans-serif';
      src: url(https://rsms.me/inter/font-files/Inter-Regular.woff2?v=3.19) format('woff2');
    }
    * {
      font-family: 'Inter', sans-serif;
    }
  </style>
</head>
<body style="margin:0">
  <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:600px;min-width:300px;width:100%;margin-left:auto;margin-right:auto;padding:0.5rem">
    <tbody>
      <tr style="width:100%">
        <td>
          <h2 style="margin:0 0 12px 0;text-align:left;color:#111827;font-size:30px;line-height:36px;font-weight:700">
            <strong>创剧星球</strong>
          </h2>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">
            欢迎您加入创剧星球,您的邮件验证码为:
`
	emailTemplateAfter := `
          </p>
          <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:100%;text-align:left;margin-bottom:0px">
            <tbody>
              <tr style="width:100%">
                <td>
                  <a href="https://ic.1024110.xyz" style="line-height:100%;text-decoration:none;display:inline-block;max-width:100%;color:#ffffff;background-color:#000000;border:2px solid #000000;font-size:14px;font-weight:500;border-radius:9999px;padding:12px 32px" target="_blank">
                    <span style="max-width:100%;display:inline-block;line-height:120%;mso-padding-alt:0px;mso-text-raise:9px">
                      前往官网 ->
                    </span>
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
          <div style="height:64px"></div>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">
            感谢您使用我们强大的由AI驱动的工程化剧集开发工具和社区。
          </p>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">
            创剧星球团队
          </p>
        </td>
      </tr>
    </tbody>
  </table>
</body>
</html>
`
	// 使用验证码替换模板中的占位符
	emailContent := emailTemplateFront + code + emailTemplateAfter

	if err := utility.SendEmail(reqBody.Email, "创剧星球验证码", emailContent); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "邮件系统发送验证码时发生错误"))
		return
	}
	c.JSON(200, entity.SuccessResponse("验证码已发送至您的邮箱，请前往查看～"))
}
func Register(c *gin.Context) {
	var registerBody RegisterRequestBody
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "解构JSON时发生错误，请检查您的表单格式"))
		return
	}
	isHasExisted, err := CheckIfCodeExists(registerBody.Email)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if !isHasExisted {
		c.JSON(200, entity.ErrorResponse[string](400, "你需要先请求您的邮箱验证码"))
		return
	}
	code, err := GetCodeFromRedis(registerBody.Email)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if registerBody.Code != code {
		c.JSON(200, entity.ErrorResponse[string](400, "您提交的验证码与邮件验证码不符"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "加密用户密码时发生错误，请稍后重试。"))
		return
	}

	user := entity.User{
		Username:   registerBody.Username,
		Password:   string(hashedPassword),
		Email:      registerBody.Email,
		Tokens:     0,
		Permission: 0,
		Group:      0,
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Where("username = ?", user.Username).First(&user).Error; err == nil {
		tx.Rollback()
		c.JSON(200, entity.ErrorResponse[string](400, "用户名已存在，更换一个吧～"))
		return
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(200, entity.ErrorResponse[string](400, "创建用户时发生错误，请稍后重试。详细信息："+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(200, entity.ErrorResponse[string](400, "创建用户时发生错误，请稍后重试。详细信息："+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse("好极了！用户创建成功，欢迎您来到创剧星球!"))
	if err := DeleteCodeToRedis(registerBody.Email); err != nil {
		fmt.Println("注销验证码时发生错误:" + err.Error())
	}
}
