package v1

import (
	"encoding/csv"
	_ "gin_blog/docs"
	"gin_blog/pkg/app"
	"gin_blog/pkg/e"
	"gin_blog/pkg/export"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/setting"
	"gin_blog/pkg/util"
	"gin_blog/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// @Summary 获取符合条件的标签列表
// @tags tag
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{"lists":[{"id":3,"created_at":1516849721,"modified_at":0,"name":"3333","created_by":"4555","modified_by":"","state":0}],"total":29},"message":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{c}
	tagService := tag_service.Tag{}
	valid := validation.Validation{}

	tagService.PageNum = util.GetPage(c)
	tagService.PageSize = setting.AppSetting.PageSize

	name := c.Query("name")
	if name != "" {
		tagService.Name = name
		valid.MaxSize(name, 100, "name").Message("标签名最大100字符")
	}

	tagService.State = -1
	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		tagService.State = state
		valid.Range(state, 0, 1, "state").Message("状态只能0或者1")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	data := make(map[string]interface{})
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	total, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}
	data["lists"] = tags
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary 新增文章标签
// @tags tag
// @Produce json
// @Param name PostForm string true "Name"
// @Param state PostForm int false "State"
// @Param created_by PostForm int false "CreatedBy"
// @Success 200 {string} json "{"code":200, "data":{}, "message":"ok""
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	appG := app.Gin{c}
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")
	valid.Required(createdBy, "createdBy").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "createdBy").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}

	exist, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改文章标签
// @tags tag
// @Produce json
// @Param id path int true "ID"
// @Param name formData string true "Name"
// @Param modified_by formData string true "modified_by"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	appG := app.Gin{c}
	//id从url路径中获取
	id := com.StrTo(c.Param("id")).MustInt()
	//其他参数从post中获取
	modifiedBy := c.PostForm("modified_by")
	name := c.PostForm("name")

	valid := validation.Validation{}
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")
	}
	valid.Min(id, 1, "id").Message("ID不能为空且必须是大于0的整数")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         id,
		Name:       name,
		State:      state,
		ModifiedBy: modifiedBy,
	}

	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除文章标签
// @tags tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID不能为空且是必须大于0的整数")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}

	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func MakeCsv() {
	//新建csv文件
	f, err := os.Create(export.GetExcelFullPath() + "test.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//标识utf8编码格式，否则汉子会乱码
	f.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(f)
	data := [][]string{
		{"1", "test1", "test1-1"},
		{"2", "test2", "test2-1"},
		{"3", "test3", "test3-1"},
		{"4", "test4", "test4-1"},
	}

	//向csv写入内容
	w.WriteAll(data)
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		logging.Info("tag export error: ", err)
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
