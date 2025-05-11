package models

import (
	"bubble/dao"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
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

	//fmt.Println("Query result:")
	//for _, world := range worldList {
	//	fmt.Printf("World: %+v\n", world)
	//}

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

func ConfirmTemplate(tid int, wid int) (err error) {

	var template Background
	var world World

	if err = dao.DB.Where("bid = ?", tid).First(&template).Error; err != nil {
		return err
	}
	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	world.WPicture = template.BPicture

	if err = dao.DB.Save(&world).Error; err != nil {
		return err
	}

	return nil
}

func ConfirmPicture(background string, wid int) (err error) {
	var world World

	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	world.WPicture = json.RawMessage(background)

	if err = dao.DB.Save(&world).Error; err != nil {
		return err
	}

	return nil
}

func GetThisWorld(wid int) (world World, err error) {
	if err = dao.DB.Where("wid = ?", wid).Find(&world).Error; err != nil {
		return World{}, err
	}
	return
}

func GetPicture(wid int) ([][]string, error) {
	var world World
	if err := dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	var wpicture [][]string
	if err := json.Unmarshal(world.WPicture, &wpicture); err != nil {
		return nil, err
	}

	return wpicture, nil
}
