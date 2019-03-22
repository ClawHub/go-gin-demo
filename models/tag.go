package models

import (
	"github.com/jinzhu/gorm"
	"go-gin-demo/pkg/gmysql"
)

type Tag struct {
	gmysql.Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

//ByName判断标签是否存在
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := gmysql.DB.Select("id").Where("name = ? AND deleted_on = ? ", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

//增加标签
func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	if err := gmysql.DB.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

//获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = gmysql.DB.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = gmysql.DB.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

//获取标签数量
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := gmysql.DB.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

//根据id判断标签是否存在
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := gmysql.DB.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

//删除标签
func DeleteTag(id int) error {
	if err := gmysql.DB.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}

//编辑标签
func EditTag(id int, data interface{}) error {
	if err := gmysql.DB.Model(&Tag{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

//清空所有标签
func CleanAllTag() (bool, error) {
	if err := gmysql.DB.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}

	return true, nil
}
