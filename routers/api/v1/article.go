package v1

import (
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

//获取单个文章【指定文章】
func GetArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()

	valid.Min(id, 1, "id").Message("ID必须是大于0的整数")
	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			data = models.GetArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Get article:", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

//获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var tagId int = -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须是大于0的整数")
	}

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只能是0或者1")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logging.Info("Get articles:", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})

}

//新增文章
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")

	valid.Min(tagId, 1, "tag_id").Message("标签ID必须是大于0的整数")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")

	valid.MaxSize(title, 100, "title").Message("标题最长100字符")
	valid.MaxSize(desc, 255, "desc").Message("标题最长255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长65535字符")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			code = e.SUCCESS
			models.AddArticle(data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info("Add Article", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	modifiedBy := c.PostForm("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(c.PostForm("state")).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")
	}
	valid.Min(id, 1, "id").Message("ID必须是正整数")
	valid.MaxSize(title, 100, "title").Message("标题最长100字符")
	valid.MaxSize(desc, 255, "desc").Message("标题最长255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长255字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			if tagId > 0 {
				data["tag_id"] = tagId
			}
			if title != "" {
				data["title"] = title

			}
			if desc != "" {
				data["desc"] = desc

			}
			if content != "" {
				data["content"] = content

			}
			if state != -1 {
				data["state"] = state

			}
			data["modified_by"] = modifiedBy

			code = e.SUCCESS
			models.EditArticle(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info("Edit Article", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})

}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须是正整数")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			models.DeleteArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("Delete article:", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}
