package article_service

import (
	"encoding/json"
	"gin_blog/models"
	"gin_blog/pkg/gredis"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/setting"
	"gin_blog/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	data := make(map[string]interface{})
	if a.TagID > 0 {
		data["tag_id"] = a.TagID
	}

	if a.Title != "" {
		data["title"] = a.Title
	}

	if a.Desc != "" {
		data["desc"] = a.Desc
	}

	if a.Content != "" {
		data["content"] = a.Content
	}

	if a.State != -1 {
		data["state"] = a.State
	}

	if a.CoverImageUrl != "" {
		data["cover_image_url"] = a.CoverImageUrl
	}
	return models.EditArticle(a.ID, data)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	//先从redis中查找
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}
	//redis无数据 从数据库查找并放入redis
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	//将从数据库中拿的数据放入redis
	gredis.Set(key, article, setting.RedisSetting.ExpireTime)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	keys := cache.GetArticlesKey()
	//从redis取数据,不存在则从数据库获取
	if gredis.Exists(keys) {
		data, err := gredis.Get(keys)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	//将数据放入redis
	gredis.Set(keys, articles, setting.RedisSetting.ExpireTime)
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
