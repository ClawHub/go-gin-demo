package auth_service

import "go-gin-demo/models"

//此类型用于业务传参
type Auth struct {
	Username string
	Password string
}

//鉴权
func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}
