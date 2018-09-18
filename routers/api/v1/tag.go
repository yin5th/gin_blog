package v1

import (
	_ "gin_blog/docs"
	"gin_blog/models"
	"gin_blog/pkg/e"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/setting"
	"gin_blog/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

//获取文章标签列表
func GetTags(c *gin.Context) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	name := c.Query("name")
	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("statue"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

// @Summary 新增文章标签
// @Produce json
// @Param name PostForm string true "Name"
// @Param state PostForm int false "State"
// @Param created_by PostForm int false "CreatedBy"
// @Success 200 {string} json "{"code":200, "data":{}, "message":"ok""
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")
	valid.Required(createdBy, "createdBy").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "createdBy").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Add Tag: ", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}

// @Summary 修改文章标签
// @Produce json
// @Param id path int true "ID"
// @Param name formData string true "Name"
// @Param modified_by formData string true "modified_by"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	modifiedBy := c.PostForm("modified_by")
	name := c.PostForm("name")

	valid := validation.Validation{}
	var state = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")
	}
	valid.Min(id, 1, "id").Message("ID不能为空且必须是大于0的整数")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			if state != -1 {
				data["state"] = state
			}
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Edit Tag: ", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID不能为空且是必须大于0的整数")
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Delete Tag: ", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}
