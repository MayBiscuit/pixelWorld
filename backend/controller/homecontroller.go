package controller

import (
	"bubble/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllWorldList(c *gin.Context) {
	id, ok := c.Params.Get("userid")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	worldList, err := models.GetAllWorld(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, worldList)
	}
}

func GetIngWorldList(c *gin.Context) {
	id, ok := c.Params.Get("userid")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	worldList, err := models.GetIngWorld(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, worldList)
	}
}

func GetEdWorldList(c *gin.Context) {
	id, ok := c.Params.Get("userid")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	worldList, err := models.GetEdWorld(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, worldList)
	}
}

func GetSearchWorldList(context *gin.Context) {
	id, ok := context.Params.Get("userid")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	key := context.Query("key")
	if key == "" {
		context.JSON(http.StatusOK, gin.H{"error": "无效的key"})
		return
	}

	worldList, err := models.GetSearchWorldList(id, key)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, worldList)
	}
}
