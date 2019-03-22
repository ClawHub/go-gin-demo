package api

import (
	"go-gin-demo/pkg/app"
	"go-gin-demo/pkg/e"
	"go-gin-demo/pkg/util"
	"go-gin-demo/service/auth_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

//此结构用于参数校验
type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

//鉴权
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	//validation is a form validation for a data validation and error collecting using Go.
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		//记录错误信息
		app.MarkErrors(valid.Errors)
		//返回参数校验失败
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//鉴权服务
	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	//生成token,生校时间3小时
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
