package api

import (
	"fmt"
	"gin_blog/models"
	"gin_blog/pkg/e"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid := validation.Validation{}
	//验证用户和密码是否符合条件
	confidition := Auth{Username: username, Password: password}
	ok, _ := valid.Valid(confidition)

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if ok {
		isExist := models.CheckAuth(username, password)
		fmt.Println(isExist)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			//auth账号或密码错误
			code = e.ERROR_AUTH_FAIL
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Auth token:", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}
