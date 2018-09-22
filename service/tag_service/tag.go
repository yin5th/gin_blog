package tag_service

import (
	"encoding/json"
	"gin_blog/models"
	"gin_blog/pkg/export"
	"gin_blog/pkg/gredis"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/setting"
	"gin_blog/service/cache_service"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/henrylee2cn/pholcus/common/xlsx"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	tag := make(map[string]interface{})
	tag["modified_by"] = t.ModifiedBy
	if t.State >= 0 {
		tag["state"] = t.State
	}
	if t.Name != "" {
		tag["name"] = t.Name
	}

	return models.EditTag(t.ID, tag)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		Name:     t.Name,
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()

	//查看redis是否有数据
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	//redis无数据 从数据库获取
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, setting.RedisSetting.ExpireTime)
	return tags, nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedAt),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedAt),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	timeNow := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + timeNow + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {
	xlsxFile, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsxFile.GetRows("标签信息")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}

	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
