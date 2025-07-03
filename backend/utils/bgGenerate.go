package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type APIConfig struct {
	GenerateURL string
	APIKey      string
}

type GenerateRequest struct {
	Model      string                 `json:"model"`
	Input      map[string]string      `json:"input"`
	Parameters map[string]interface{} `json:"parameters"`
}

func averageColor(img *image.RGBA, rect image.Rectangle) color.RGBA {
	var r, g, b, a, count uint32

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			r += uint32(c.R)
			g += uint32(c.G)
			b += uint32(c.B)
			a += uint32(c.A)
			count++
		}
	}

	if count > 0 {
		r /= count
		g /= count
		b /= count
		a /= count
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func toHexColor(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func GenerateByPicture(img *image.RGBA, gridSize int) ([][]string, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	pixelated := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(pixelated, pixelated.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// 将图像裁剪成正方形
	sideLength := min(width, height)
	cropped := image.NewRGBA(image.Rect(0, 0, sideLength, sideLength))
	draw.Draw(cropped, cropped.Bounds(), img, image.Point{X: (width - sideLength) / 2, Y: (height - sideLength) / 2}, draw.Src)

	// 重新计算 gridSize
	gridSize = sideLength / gridSize

	var colors [][]string
	for y := 0; y < height; y += gridSize {
		var row []string
		for x := 0; x < width; x += gridSize {
			rect := image.Rect(x, y, min(x+gridSize, width), min(y+gridSize, height))
			avgColor := averageColor(img, rect)
			row = append(row, toHexColor(avgColor))

			for py := rect.Min.Y; py < rect.Max.Y; py++ {
				for px := rect.Min.X; px < rect.Max.X; px++ {
					pixelated.Set(px, py, avgColor)
				}
			}
		}
		colors = append(colors, row)
	}

	return colors, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GenerateByDescription(description string, gridSize int) ([][]string, error) {
	fmt.Println("开始调用文生图")
	config := APIConfig{
		GenerateURL: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis",
		APIKey:      "", // APIKEY有效期够用到暑假结束
	}

	taskID, err := sendGenerateRequest(config, description)
	fmt.Println("传回ai的taskID:", taskID)
	if err != nil {
		return nil, err
	}

	imgUrl, _ := getImageUrl(config, taskID, 5)
	fmt.Println("传回imgUrl:", imgUrl)

	img, _ := getImage(imgUrl)
	fmt.Println("传回img:", img)

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(img.Bounds())
		draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)
	}

	background, _ := GenerateByPicture(rgbaImg, gridSize)
	fmt.Println("生成background:", background)

	return background, nil
}

func sendGenerateRequest(config APIConfig, description string) (string, error) {

	requestBody := map[string]interface{}{
		"model": "wanx2.1-t2i-turbo",
		"input": map[string]string{
			"prompt": description,
		},
		"parameters": map[string]interface{}{
			"size": "1024*1024",
			"n":    1,
		},
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal("Error marshalling request body:", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", config.GenerateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// 设置请求头
	req.Header.Set("X-DashScope-Async", "enable")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Received non-200 response status:", resp.StatusCode)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// 解析JSON响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	// 打印结果
	fmt.Printf("Response: %+v\n", result)

	// 提取任务ID
	output, ok := result["output"].(map[string]interface{})
	if !ok {
		log.Fatal("Error asserting output type")
	}
	taskID, ok := output["task_id"].(string)
	if !ok {
		log.Fatal("Error asserting taskid type")
	}

	return taskID, nil
}

func getImageUrl(config APIConfig, taskID string, retries int) (string, error) {
	fmt.Println("getImageUrl被调用")
	// 构建查询任务的URL
	urlQuery := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID)

	// 创建请求
	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	// 设置请求头
	req.Header.Add("Authorization", "Bearer "+config.APIKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received non-200 response status: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	// 解析JSON响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("Error decoding JSON: %v", err)
	}

	// 检查任务状态
	status, ok := result["output"].(map[string]interface{})["task_status"].(string)
	if !ok {
		return "", fmt.Errorf("Error asserting task status type")
	}

	// 如果任务未完成，且还有重试次数，则等待1秒后重试
	if status == "RUNNING" || status != "SUCCEEDED" && retries > 0 {
		time.Sleep(1 * time.Second)
		return getImageUrl(config, taskID, retries-1)
	}

	// 如果任务未完成且没有重试次数，返回错误
	if status != "SUCCEEDED" {
		return "", fmt.Errorf("Task status is %s, no retries left", status)
	}

	// 打印结果
	//fmt.Println("1 ", result)
	output := result["output"].(map[string]interface{})
	results := output["results"].([]interface{})
	//fmt.Println("2 ", results)
	firstResult := results[0].(map[string]interface{})
	//fmt.Println("3 ", firstResult)
	urlImage := firstResult["url"].(string)
	//fmt.Println("4 ", urlImage)
	return urlImage, nil
}

func getImage(url string) (image.Image, error) {
	// 发送HTTP GET请求获取图像
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("无法下载图像: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close() // 确保响应体在函数退出时关闭

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("下载图像失败，状态码: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("下载图像失败，状态码: %d", resp.StatusCode)
	}

	// 解码图像
	img, format, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Printf("图像解码失败: %v\n", err)
		return nil, err
	}

	// 输出图像的格式和尺寸
	fmt.Printf("图像格式: %s\n", format)
	fmt.Printf("图像尺寸: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

	return img, nil
}
