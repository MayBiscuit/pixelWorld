package controller

import (
	"bubble/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateWorld(c *gin.Context) {

	type CreateWorldRequest struct {
		UserID string `json:"userid" binding:"required"`
		WName  string `json:"wname" binding:"required"`
		WDesc  string `json:"wdesc" binding:"required"`
	}

	var req CreateWorldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newWorld := models.World{
		UserID:  req.UserID,
		WName:   req.WName,
		WDesc:   req.WDesc,
		WStatus: "绘制中",
		WSize:   25,
	}

	if err := models.CreateWorld(&newWorld); err != nil {
		// 如果保存失败，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "World created successfully", "world": newWorld})

}

func GetAllTemplateList(context *gin.Context) {

}

func GetAllTypeTemplateList(context *gin.Context) {

}

func GetSearchTemplateList(context *gin.Context) {

}

func GetSearchTypeTemplateList(context *gin.Context) {

}

func GetChooseTemplate(context *gin.Context) {

}

func ConfirmTemplate(context *gin.Context) {

}

func UploadPicture(context *gin.Context) {

}

func ConfirmPicture(context *gin.Context) {

}

func UploadDescription(context *gin.Context) {

}

func ConfirmDescription(context *gin.Context) {

}
