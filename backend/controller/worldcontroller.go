package controller

import (
	"bubble/models"
	"bubble/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
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

func GetAllTemplateList(c *gin.Context) {
	templateList, err := models.GetAllTemplate()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, templateList)
	}
}

func GetAllTypeTemplateList(c *gin.Context) {
	category := c.Query("type")

	templateList, err := models.GetAllTypeTemplate(category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, templateList)
	}
}

func GetSearchTemplateList(c *gin.Context) {
	key := c.Query("key")

	templateList, err := models.GetSearchTemplateList(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, templateList)
	}
}

func GetSearchTypeTemplateList(c *gin.Context) {
	key := c.Query("key")
	category := c.Query("type")

	templateList, err := models.GetSearchTypeTemplateList(key, category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, templateList)
	}
}

func GetChooseTemplate(c *gin.Context) {
	tidStr := c.Query("bid")
	tid, _ := strconv.Atoi(tidStr)

	background, err := models.GetChooseTemplate(tid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, background)
	}
}

func ConfirmTemplate(c *gin.Context) {
	tidStr := c.Query("tid")
	tid, _ := strconv.Atoi(tidStr)
	widStr := c.Query("wid")
	wid, _ := strconv.Atoi(widStr)

	err := models.ConfirmTemplate(tid, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"meg": "confirm success"})
	}
}

func UploadPicture(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件上传失败",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法打开文件",
		})
		return
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图像解码失败: " + err.Error()})
		return
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(img.Bounds())
		draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)
	}

	background, _ := utils.GenerateByPicture(rgbaImg, 10)

	c.JSON(http.StatusOK, background)

}

func ConfirmPicture(c *gin.Context) {
	background := c.Query("background")
	widStr := c.Query("wid")
	wid, _ := strconv.Atoi(widStr)

	err := models.ConfirmPicture(background, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"meg": "confirm success"})
	}
}

func UploadDescription(c *gin.Context) {
	description := c.Query("description")
	fmt.Fprint(c.Writer, description)
	background, _ := utils.GenerateByDescription(description)

	c.JSON(http.StatusOK, background)
}

func ConfirmDescription(c *gin.Context) {
	background := c.Query("background")
	widStr := c.Query("wid")
	wid, _ := strconv.Atoi(widStr)

	err := models.ConfirmPicture(background, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"meg": "confirm success"})
	}
}
