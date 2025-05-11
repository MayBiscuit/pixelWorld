package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

//func GenerateByPicture(file *multipart.FileHeader) interface{} {
//	// TODO 图像素化
//	return nil
//}

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

	// Create a new image for the pixelated version
	pixelated := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(pixelated, pixelated.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

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

func GenerateByDescription(description string) ([][]string, error) {
	config := APIConfig{
		GenerateURL: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis",
		APIKey:      "",
	}

	requestBody := GenerateRequest{
		Model: "wanx2.1-t2i-turbo",
		Input: map[string]string{
			"prompt": description,
		},
		Parameters: map[string]interface{}{},
	}

	taskID, err := sendGenerateRequest(config, requestBody)
	fmt.Println(taskID)
	if err != nil {
		return nil, err
	}

	return getImage(config, taskID, 5)

}

func sendGenerateRequest(config APIConfig, requestBody GenerateRequest) (string, error) {
	bodyBytes, err := json.Marshal(requestBody)
	fmt.Println("a" + string(bodyBytes))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.GenerateURL, bytes.NewBuffer(bodyBytes))
	fmt.Println(req.Header)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("generate request failed with status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("JSON解码失败: %v\n", err)
		return "", err
	}
	fmt.Println(result)

	taskID, ok := result["task_id"].(string)
	if !ok {
		return "", fmt.Errorf("task ID not found in response")
	}

	return taskID, nil
}

func getImage(config APIConfig, taskID string, retries int) ([][]string, error) {
	getImageURL := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis?task_id=%s", taskID)

	for i := 0; i < retries; i++ {
		req, err := http.NewRequest("GET", getImageURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+config.APIKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, err
			}

			status, ok := result["status"].(string)
			if !ok || status == "RUNNING" || status != "SUCCEEDED" {
				if i < retries-1 {
					time.Sleep(1 * time.Second)
					continue
				}
				return nil, fmt.Errorf("image generation failed with status: %s", status)
			}

			imageURL := result["url"].(string)

			// 发送HTTP GET请求获取图像
			resp, err := http.Get(imageURL)
			if err != nil {
				fmt.Printf("无法下载图像: %v\n", err)
				return nil, nil
			}
			defer resp.Body.Close() // 确保响应体在函数退出时关闭

			// 检查HTTP响应状态码
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("下载图像失败，状态码: %d\n", resp.StatusCode)
				return nil, nil
			}

			// 解码图像
			img, format, err := image.Decode(resp.Body)
			if err != nil {
				fmt.Printf("图像解码失败: %v\n", err)
				return nil, nil
			}

			// 输出图像的格式和尺寸
			fmt.Printf("图像格式: %s\n", format)
			fmt.Printf("图像尺寸: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

			rgbaImg, ok := img.(*image.RGBA)
			if !ok {
				rgbaImg = image.NewRGBA(img.Bounds())
				draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)
			}

			background, _ := GenerateByPicture(rgbaImg, 10)

			return background, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("failed to retrieve image after %d attempts", retries)
}
