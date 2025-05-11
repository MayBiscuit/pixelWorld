package models

import (
	"bubble/dao"
	"encoding/json"
)

type Sticker struct {
	SID      int             `json:"sid" gorm:"column:sid;primary_key;AUTO_INCREMENT"`
	SName    string          `json:"sname" gorm:"column:sname"`
	SType    string          `json:"stype" gorm:"column:stype"`
	SPicture json.RawMessage `json:"spicture" gorm:"column:spicture"`
}

func GetAllSticker() (stickerList []Background, err error) {
	if err = dao.DB.Find(&stickerList).Error; err != nil {
		return nil, err
	}

	//fmt.Println("Query result:")
	//for _, world := range StickerList {
	//	fmt.Printf("StickerList: %+v\n", world)
	//}

	return
}

func GetAllTypeSticker(category string) (stickerList []Background, err error) {
	if err = dao.DB.Where("btype = ?", category).Find(&stickerList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchStickerList(key string) (stickerList []Background, err error) {
	if err = dao.DB.Where("bname LIKE ?", "%"+key+"%").Find(&stickerList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchTypeStickerList(key string, category string) (stickerList []Background, err error) {
	if err = dao.DB.Where("btype = ?", category).Where("bname LIKE ?", "%"+key+"%").Find(&stickerList).Error; err != nil {
		return nil, err
	}
	return
}

func GetChooseSticker(tid int) (background Background, err error) {
	if err = dao.DB.Where("bid = ?", tid).Find(&background).Error; err != nil {
		return Background{}, err
	}
	return
}
