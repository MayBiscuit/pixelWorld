package models

import (
	"bubble/dao"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
)

type Background struct {
	BID      int             `json:"bid" gorm:"column:bid;primary_key;AUTO_INCREMENT"`
	BName    string          `json:"bname" gorm:"column:bname"`
	BType    string          `json:"btype" gorm:"column:btype"`
	BPicture json.RawMessage `json:"bpicture" gorm:"column:bpicture"`
}

func GetAllTemplate() (templateList []Background, err error) {
	if err = dao.DB.Find(&templateList).Error; err != nil {
		return nil, err
	}

	//fmt.Println("Query result:")
	//for _, world := range templateList {
	//	fmt.Printf("TemplateList: %+v\n", world)
	//}

	return
}

func GetAllTypeTemplate(category string) (templateList []Background, err error) {
	if err = dao.DB.Where("btype = ?", category).Find(&templateList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchTemplateList(key string) (templateList []Background, err error) {
	if err = dao.DB.Where("bname LIKE ?", "%"+key+"%").Find(&templateList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchTypeTemplateList(key string, category string) (templateList []Background, err error) {
	if err = dao.DB.Where("btype = ?", category).Where("bname LIKE ?", "%"+key+"%").Find(&templateList).Error; err != nil {
		return nil, err
	}
	return
}

func GetChooseTemplate(tid int) (background Background, err error) {
	if err = dao.DB.Where("bid = ?", tid).Find(&background).Error; err != nil {
		return Background{}, err
	}
	return
}

func GetTemplate(tid int) ([][]string, error) {
	var bg Background
	if err := dao.DB.Where("bid = ?", tid).First(&bg).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	var bpicture [][]string
	if err := json.Unmarshal(bg.BPicture, &bpicture); err != nil {
		return nil, err
	}

	return bpicture, nil
}
