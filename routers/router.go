package routers

import (
	"gin_blog/middleware/jwt"
	"gin_blog/pkg/setting"
	"gin_blog/pkg/upload"
	"gin_blog/routers/api"
	"gin_blog/routers/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//设置模式
	gin.SetMode(setting.ServerSetting.RunMode)

	//静态资源
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	//获取token
	r.POST("/auth", api.GetAuth)
	//api文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//上传文件
	r.POST("/upload", api.UploadImage)

	//引入jwt中间件
	r.Use(jwt.JWT())
	apiv1 := r.Group("/api/v1")
	{
		/*****标签*****/
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		/*****文章*****/
		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//修改文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
