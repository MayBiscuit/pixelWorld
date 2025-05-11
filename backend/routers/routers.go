package routers

import (
	"bubble/controller"
	"bubble/setting"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/allworld/:userid", controller.GetAllWorldList)
		homeGroup.GET("/ingworld/:userid", controller.GetIngWorldList)
		homeGroup.GET("/edworld/:userid", controller.GetEdWorldList)
		homeGroup.GET("/searchworld/:userid", controller.GetSearchWorldList)
	}

	worldGroup := r.Group("/world")
	{
		worldGroup.POST("/createworld", controller.CreateWorld)
		//worldGroup.GET("/all", controller.GetAllTemplateList)

		templateGroup := worldGroup.Group("/template")
		{
			templateGroup.GET("/all", controller.GetAllTemplateList)
			templateGroup.GET("/alltype", controller.GetAllTypeTemplateList)
			templateGroup.GET("/search", controller.GetSearchTemplateList)
			templateGroup.GET("/searchtype", controller.GetSearchTypeTemplateList)
			templateGroup.GET("/choose", controller.GetChooseTemplate)
			templateGroup.PUT("/confirm", controller.ConfirmTemplate)
		}

		pictureGroup := worldGroup.Group("/picture")
		{
			pictureGroup.POST("/upload", controller.UploadPicture)
			pictureGroup.PUT("/confirm", controller.ConfirmPicture)
		}

		aiGroup := worldGroup.Group("/ai")
		{
			aiGroup.POST("/upload", controller.UploadDescription)
			aiGroup.PUT("/confirm", controller.ConfirmDescription)
		}
	}

	drawGroup := r.Group("/draw")
	{
		drawGroup.GET("/thisworld/:wid", controller.GetThisWorld)
		drawGroup.PUT("/draw", controller.Draw)

		drawGroup.GET("/colorrank", controller.GetColorRank)
		drawGroup.GET("/save", controller.SavePicture)

		stickerGroup := drawGroup.Group("/sticker")
		{
			stickerGroup.GET("/all", controller.GetAllStickerList)
			stickerGroup.GET("/alltype", controller.GetAllTypeStickerList)
			stickerGroup.GET("/search", controller.GetSearchStickerList)
			stickerGroup.GET("/searchtype", controller.GetSearchTypeStickerList)
			stickerGroup.GET("/choose", controller.GetChooseSticker)
		}
	}

	return r
}
