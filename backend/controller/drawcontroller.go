package controller

import (
	"bubble/models"
	"bubble/utils"
	"fmt"
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

//func Draw(c *gin.Context) {
//	background := c.Query("background")
//	widStr := c.Query("wid")
//	wid, _ := strconv.Atoi(widStr)
//
//	err := models.ConfirmPicture(background, wid)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
//	} else {
//		c.JSON(http.StatusOK, gin.H{"meg": "draw success"})
//	}
//}

func Draw(c *gin.Context) {
	var req struct {
		Background string `json:"background"`
		Wid        int    `json:"wid"`
		Wsize      int    `json:"wsize"`
	}

	fmt.Println("1")

	// 绑定请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	fmt.Println("2")

	background := req.Background
	wid := req.Wid
	wsize := req.Wsize

	fmt.Println("put wid: ", wid)

	err := models.ConfirmPicture(background, wid, wsize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		fmt.Println("3")
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "draw success"})
		fmt.Println("4")
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

func DeleteThisWorld(c *gin.Context) {
	widParam := c.Param("wid")
	wid, err := strconv.Atoi(widParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 wid"})
		return
	}

	// 调用模型层的函数删除世界
	err = models.DeleteWorld(wid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func ChangeWorldStatus(c *gin.Context) {
	widParam := c.Param("wid")
	wid, err := strconv.Atoi(widParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 wid"})
		return
	}

	// 调用模型层的函数更新 wstatus
	err = models.ChangeWorldStatus(wid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"msg": "状态更新成功"})
}

func ModifyWorldName(c *gin.Context) {
	var req struct {
		Wname string `json:"wname"`
		Wid   int    `json:"wid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	wname := req.Wname
	wid := req.Wid

	err := models.ModifyWorldName(wname, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "modify wname success"})
		fmt.Println("4")
	}
}

func ModifyWorldDesc(c *gin.Context) {
	var req struct {
		Wdesc string `json:"wdesc"`
		Wid   int    `json:"wid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	wdesc := req.Wdesc
	wid := req.Wid

	err := models.ModifyWorldDesc(wdesc, wid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "modify wdesc success"})
		fmt.Println("4")
	}
}
