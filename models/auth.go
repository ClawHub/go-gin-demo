package models

import (
	"github.com/jinzhu/gorm"
	"go-gin-demo/pkg/gmysql"
)

//此类型用于与数据库交互
type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//鉴权
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := gmysql.DB.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}
