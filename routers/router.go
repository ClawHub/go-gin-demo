package routers

import (
	"github.com/gin-gonic/gin"
	"go-gin-demo/middleware/cors"
	"go-gin-demo/middleware/jwt"
	"go-gin-demo/pkg/export"
	"go-gin-demo/pkg/setting"
	"go-gin-demo/pkg/upload"
	"go-gin-demo/routers/api"
	"go-gin-demo/routers/api/v1"
	"net/http"
)

//初始化路由
func InitRouter() *gin.Engine {

	//默认初始化 Gin
	r := gin.New()
	//Logger实例将日志写入gin.DefaultWriter的日志记录器中间件。
	r.Use(gin.Logger())

	//Recovery返回一个中间件，该中间件从任何恐慌中恢复，如果有500，则写入500。
	r.Use(gin.Recovery())
	// 使用跨域中间件
	r.Use(cors.Cors())
	//设置mode-----"debug","release","test"
	gin.SetMode(setting.ServerSetting.RunMode)
	//工程名
	project := r.Group(setting.ServerSetting.ProjectName)
	//静态文件服务
	project.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	project.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	//健康检查
	project.GET("/welcome", api.Welcome)
	//鉴权
	project.GET("/auth", api.GetAuth)
	//上传图片
	project.POST("/upload", api.UploadImage)
	//路由组
	apiv1 := project.Group("/api/v1")

	//使用JSON Web Tokens 中间件
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
