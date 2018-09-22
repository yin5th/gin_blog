package api

import (
	"gin_blog/models"
	"gin_blog/pkg/app"
	"gin_blog/pkg/e"
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
	appG := app.Gin{c}
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid := validation.Validation{}
	//验证用户和密码是否符合条件
	confidition := Auth{Username: username, Password: password}
	ok, _ := valid.Valid(confidition)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	data := make(map[string]interface{})

	isExist, err := models.CheckAuth(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	data["token"] = token
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
