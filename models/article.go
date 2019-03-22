package models

import (
	"github.com/jinzhu/gorm"
	"go-gin-demo/pkg/gmysql"
)

type Article struct {
	gmysql.Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

//根据ID判断文章是否存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := gmysql.DB.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

//获取文章数量
func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := gmysql.DB.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

//获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := gmysql.DB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

//获取指定文章
func GetArticle(id int) (*Article, error) {
	var article Article
	err := gmysql.DB.Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

//更新指定文章
func EditArticle(id int, data interface{}) error {
	if err := gmysql.DB.Model(&Article{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

//新建文章
func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}
	if err := gmysql.DB.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

//删除指定文章
func DeleteArticle(id int) error {
	if err := gmysql.DB.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}

	return nil
}

//清空所有文章
func CleanAllArticle() error {
	if err := gmysql.DB.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{}).Error; err != nil {
		return err
	}

	return nil
}
