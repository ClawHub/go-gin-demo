package article_service

import (
	"encoding/json"
	"go-gin-demo/models"
	"go-gin-demo/pkg/gredis"
	"go-gin-demo/pkg/logging"
	"go-gin-demo/service/cache_service"
	"go.uber.org/zap"
)

//文章信息
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

//新建文章
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

//更新指定文章
func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	})
}

//获取文章
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	//从缓存获取
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key).Val() > 0 {
		data, err := gredis.Get(key).Bytes()
		if err != nil {
			logging.HTTPLogger.Warn("Get Article cache fail", zap.Error(err))
		} else {
			_ = json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	//从库中获取文章
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	//入缓存
	gredis.Set(key, article, 3600)
	return article, nil
}

//获取文章列表
func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	//从缓存中获取文章
	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key).Val() > 0 {
		data, err := gredis.Get(key).Bytes()
		if err != nil {
			logging.HTTPLogger.Info("cache get articles fail", zap.Error(err))
		} else {
			_ = json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	//从库中获取文章
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	//入缓存
	gredis.Set(key, articles, 3600)
	return articles, nil
}

//删除指定文章
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

//根据ID判断文章是否存在
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

//获取文章数量
func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

//组装
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
