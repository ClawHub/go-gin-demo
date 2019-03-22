package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go-gin-demo/pkg/e"
	"net/http"
)

//绑定和校验
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	// Bind检查内容类型以自动选择绑定引擎，
	//根据“Content-Type”标头使用不同的绑定:
	//“application/json”——> json绑定
	//“application/xml”——> xml绑定
	//否则——>返回一个错误。
	//如果Content-Type == "application/ JSON "使用JSON或XML作为JSON输入，则它将请求体解析为JSON。
	//它将json有效负载解码为指定为指针的结构。
	//如果输入无效，它将写入400个错误并在响应中设置Content-Type头“text/plain”。
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	//参数校验
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
