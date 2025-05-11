package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sort"
	"strconv"
)

type ColorCount struct {
	String string
	Freq   int
}

func PixelCount(arr [][]string) []ColorCount {
	freqMap := make(map[string]int)

	for _, subArr := range arr {
		for _, str := range subArr {
			freqMap[str]++
		}
	}

	freqSlice := make([]ColorCount, 0, len(freqMap))
	for str, freq := range freqMap {
		freqSlice = append(freqSlice, ColorCount{String: str, Freq: freq})
	}

	sort.Slice(freqSlice, func(i, j int) bool {
		return freqSlice[i].Freq > freqSlice[j].Freq
	})

	n := len(freqSlice)
	if n > 10 {
		freqSlice = freqSlice[:10]
	}

	return freqSlice
}

func hexToRGB(hex string) (color.RGBA, error) {
	var rgb color.RGBA
	values, err := strconv.ParseUint(hex[1:], 16, 32)
	if err != nil {
		return rgb, err
	}

	rgb.R = uint8((values >> 16) & 0xFF)
	rgb.G = uint8((values >> 8) & 0xFF)
	rgb.B = uint8(values & 0xFF)
	rgb.A = 0xFF
	return rgb, nil
}

//func SavePicture(wid int) (picturePath string, err error) {
//	// TODO 返回保存画室图片url
//	return "", nil
//}

func SavePicture(colors [][]string) ([]byte, error) {
	if len(colors) == 0 || len(colors[0]) == 0 {
		return nil, fmt.Errorf("empty color array")
	}

	width, height := len(colors[0]), len(colors)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y, row := range colors {
		for x, hex := range row {
			rgb, err := hexToRGB(hex)
			if err != nil {
				return nil, fmt.Errorf("invalid color format at (%d, %d): %v", x, y, err)
			}
			img.Set(x, y, rgb)
		}
	}
	var buf []byte
	w := bytes.NewBuffer(buf)
	err := png.Encode(w, img)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
