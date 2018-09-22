package v1

import (
	"gin_blog/pkg/app"
	"gin_blog/pkg/e"
	"gin_blog/pkg/setting"
	"gin_blog/pkg/util"
	"gin_blog/service/article_service"
	"gin_blog/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取单个文章
// @tags article
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{"id":3,"created_on":1516937037,"modified_on":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1},"message":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{c}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须是大于0的整数")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary 获取多个文章
// @tags article
// @Produce  json
// @Param tag_id query int false "TagID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":[{"id":3,"created_at":1516937037,"modified_at":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1}],"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}

	valid := validation.Validation{}

	tagId := -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须是大于0的整数")
	}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能是0或者1")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	data := map[string]interface{}{
		"lists": articles,
		"total": total,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// @Summary 新增文章
// @tags article
// @Produce  json
// @Param tag_id formData int true "TagID"
// @Param title formData string true "Title"
// @Param desc formData string true "Desc"
// @Param content formData string true "Content"
// @Param created_by formData string true "CreatedBy"
// @Param state formData int true "State"
// @Param cover_image_url formData string true "CoverImageUrl"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	coverImageUrl := c.PostForm("cover_image_url")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面地址不能为空")

	valid.Min(tagId, 1, "tag_id").Message("标签ID必须是大于0的整数")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")

	valid.MaxSize(title, 100, "title").Message("标题最长100字符")
	valid.MaxSize(desc, 255, "desc").Message("标题最长255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长65535字符")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长100字符")

	appG := app.Gin{c}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{ID: tagId}
	tagExist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !tagExist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         tagId,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CreatedBy:     createdBy,
		State:         state,
		CoverImageUrl: coverImageUrl,
	}

	err = articleService.Add()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改文章
// @tags article
// @Produce  json
// @Param id param int true "ID"
// @Param tag_id formData string false "TagID"
// @Param title formData string false "Title"
// @Param desc formData string false "Desc"
// @Param content formData string false "Content"
// @Param modified_by formData string true "ModifiedBy"
// @Param state formData int false "State"
// @Param cover_image_url formData string false "CoverImageUrl"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	appG := app.Gin{c}
	//获取传入的参数
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	coverImageUrl := c.PostForm("cover_image_url")
	modifiedBy := c.PostForm("modified_by")

	//验证传入参数合法性
	valid := validation.Validation{}

	state := -1
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
	valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("封面地址最长为255字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		ID:            id,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		ModifiedBy:    modifiedBy,
		State:         state,
	}
	//文章是否存在
	articleExist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !articleExist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	//文章标签是否存在
	if tagId > 0 {
		tagService := tag_service.Tag{ID: tagId}
		tagExist, err := tagService.ExistByID()
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
			return
		}

		if !tagExist {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
			return
		}
		articleService.TagID = tagId
	}

	//执行修改
	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除文章
// @tags article
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"message":"请求参数错误"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须是正整数")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
