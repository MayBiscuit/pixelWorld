package routers

import (
	"bubble/controller"
	"bubble/setting"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

type User struct {
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type LoginRequest struct {
	Code     string `json:"code"`
	UserInfo struct {
		NickName  string `json:"nickName"`
		AvatarUrl string `json:"avatarUrl"`
	} `json:"userInfo,omitempty"`
}

func init() {
	gob.Register(User{})
}

const (
	//appID     = "wx30effb50cf207f03"
	//appSecret = "cc751a9abb24aa37358ecf24e08648c6"
	//appID     = "wxc7616829de2cf516"
	//appSecret = "93f3ac0e1d9b3bafc36b2e28eef4c11d"
	appID     = "wx8a3e8490cf4586eb"
	appSecret = "11aad2e778a45edc15800a23515651f5"
)

func getWeChatOpenID(code string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应体失败: %v", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口返回错误: %s", result.ErrCode, result.ErrMsg)
	}

	return result.OpenID, nil
}

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 配置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://servicewechat.com")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 登录接口
	r.POST("/api/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println("请求参数错误")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    400,
					"message": "请求参数错误",
					"details": err.Error(),
				},
			})
			return
		}

		// 这里应该调用微信接口获取openid
		openid, err := getWeChatOpenID(req.Code)
		if err != nil {
			fmt.Println("微信登录失败", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    500,
					"message": "微信登录失败",
					"details": err.Error(),
				},
			})
			return
		}

		// 模拟获取openid
		//openid := fmt.Sprintf("mockopenid%d", time.Now().Unix())

		// 获取session
		session, err := store.Get(c.Request, "session-name")
		if err != nil {
			fmt.Println("获取session失败")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    500,
					"message": "获取session失败",
					"details": err.Error(),
				},
			})
			return
		}

		user := User{
			OpenID:    openid,
			Nickname:  req.UserInfo.NickName,
			AvatarURL: req.UserInfo.AvatarUrl,
		}

		// 保存用户信息到session
		session.Values["user"] = user
		if err := session.Save(c.Request, c.Writer); err != nil {
			fmt.Println("保存session失败", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    500,
					"message": "保存session失败",
					"details": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"user":    user,
		})
	})

	// 登出接口
	r.POST("/api/logout", func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session-name")
		session.Options.MaxAge = -1 // 立即删除cookie
		if err := session.Save(c.Request, c.Writer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 检查登录状态
	r.GET("/api/check-auth", func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session-name")
		user, ok := session.Values["user"].(User)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"isLoggedIn": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"isLoggedIn": true,
			"user":       user,
		})
	})

	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/allworld/:userid", controller.GetAllWorldList)
		homeGroup.GET("/ingworld/:userid", controller.GetIngWorldList)
		homeGroup.GET("/edworld/:userid", controller.GetEdWorldList)
		homeGroup.GET("/searchworld/:userid", controller.GetSearchWorldList)
		homeGroup.GET("/searchingworld/:userid", controller.GetSearchIngWorldList)
		homeGroup.GET("/searchedworld/:userid", controller.GetSearchEdWorldList)
	}

	worldGroup := r.Group("/world")
	{
		worldGroup.POST("/createworld", controller.CreateWorld)
		worldGroup.PUT("/confirmEmpty", controller.ConfirmEmpty)
		//worldGroup.GET("/all", controller.GetAllTemplateList)

		templateGroup := worldGroup.Group("/template")
		{
			templateGroup.GET("/all", controller.GetAllTemplateList)
			templateGroup.GET("/alltype", controller.GetAllTypeTemplateList)
			templateGroup.GET("/search", controller.GetSearchTemplateList)
			templateGroup.GET("/searchtype", controller.GetSearchTypeTemplateList)
			templateGroup.GET("/choose", controller.GetChooseTemplate)
			templateGroup.PUT("/confirm", controller.ConfirmTemplate)
			templateGroup.GET("/save", controller.SaveTemplate)
		}

		pictureGroup := worldGroup.Group("/picture")
		{
			pictureGroup.POST("/upload", controller.UploadPicture)
			pictureGroup.PUT("/confirm", controller.ConfirmPicture)
		}

		aiGroup := worldGroup.Group("/ai")
		{
			aiGroup.GET("/upload", controller.UploadDescription)
			aiGroup.PUT("/confirm", controller.ConfirmDescription)
		}
	}

	drawGroup := r.Group("/draw")
	{
		drawGroup.GET("/thisworld/:wid", controller.GetThisWorld)
		drawGroup.PUT("/draw", controller.Draw)
		drawGroup.PUT("/modifyname", controller.ModifyWorldName)
		drawGroup.PUT("/modifydesc", controller.ModifyWorldDesc)
		drawGroup.PUT("changeworldstatus/:wid", controller.ChangeWorldStatus)
		drawGroup.DELETE("/deleteworld/:wid", controller.DeleteThisWorld)

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
