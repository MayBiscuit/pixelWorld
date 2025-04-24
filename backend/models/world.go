package models

import (
	"bubble/dao"
	"encoding/json"
	"fmt"
)

type World struct {
	WID      int             `json:"wid" gorm:"column:wid;primary_key;AUTO_INCREMENT"`
	WName    string          `json:"wname" gorm:"column:wname"`
	WDesc    string          `json:"wdesc" gorm:"column:wdesc"`
	WStatus  string          `json:"wstatus" gorm:"column:wstatus"`
	WSize    int             `json:"wsize" gorm:"column:wsize"`
	WPicture json.RawMessage `json:"wpicture" gorm:"column:wpicture;type:json"`
	UserID   string          `json:"userid" gorm:"column:userid"`
}

func GetAllWorld(id string) (worldList []World, err error) {
	if err = dao.DB.Where("userid = ?", id).Find(&worldList).Error; err != nil {
		return nil, err
	}

	fmt.Println("Query result:")
	for _, world := range worldList {
		fmt.Printf("World: %+v\n", world)
	}

	return
}

func GetIngWorld(id string) (worldList []*World, err error) {
	if err = dao.DB.Where("wstatus = ? AND userid = ?", "绘制中", id).Find(&worldList).Error; err != nil {
		return nil, err
	}
	return
}

func GetEdWorld(id string) (worldList []*World, err error) {
	if err = dao.DB.Where("wstatus = ? AND userid = ?", "绘制完毕", id).Find(&worldList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchWorldList(id string, key string) (worldList []*World, err error) {
	if err := dao.DB.Where("userid = ?", id).Where("wdesc LIKE ? OR wname LIKE ?", "%"+key+"%", "%"+key+"%").Find(&worldList).Error; err != nil {
		return nil, err
	}
	return
}

func CreateWorld(w *World) (err error) {
	err = dao.DB.Create(&w).Error
	return
}
