package controller

import (
	"bubble/models"
	"bubble/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetThisWorld(c *gin.Context) {
	widStr, ok := c.Params.Get("wid")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	wid, _ := strconv.Atoi(widStr)

	worldList, err := models.GetThisWorld(wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, worldList)
	}
}

func Draw(c *gin.Context) {
	background := c.Query("background")
	widStr := c.Query("wid")
	wid, _ := strconv.Atoi(widStr)

	err := models.ConfirmPicture(background, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"meg": "draw success"})
	}
}

func GetColorRank(c *gin.Context) {
	widStr := c.Query("wid")

	wid, _ := strconv.Atoi(widStr)

	picture, _ := models.GetPicture(wid)

	colorCounts := utils.PixelCount(picture)

	c.JSON(http.StatusOK, colorCounts)
}

func SavePicture(c *gin.Context) {
	widStr := c.Query("wid")

	wid, _ := strconv.Atoi(widStr)

	picture, _ := models.GetPicture(wid)

	pictureBytes, err := utils.SavePicture(picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(pictureBytes)
}

func GetAllStickerList(c *gin.Context) {
	stickerList, err := models.GetAllSticker()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stickerList)
	}
}

func GetAllTypeStickerList(c *gin.Context) {
	category := c.Query("type")

	stickerList, err := models.GetAllTypeSticker(category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stickerList)
	}
}

func GetSearchStickerList(c *gin.Context) {
	key := c.Query("key")

	stickerList, err := models.GetSearchStickerList(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stickerList)
	}
}

func GetSearchTypeStickerList(c *gin.Context) {
	key := c.Query("key")
	category := c.Query("type")

	stickerList, err := models.GetSearchTypeStickerList(key, category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stickerList)
	}
}

func GetChooseSticker(c *gin.Context) {
	tidStr := c.Query("sid")
	tid, _ := strconv.Atoi(tidStr)

	background, err := models.GetChooseSticker(tid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, background)
	}
}
