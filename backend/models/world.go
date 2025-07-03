package models

import (
	"bubble/dao"
	"encoding/json"
	"errors"
	"fmt"
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
	if err = dao.DB.Where("userid = ?", id).Order("wid DESC").Find(&worldList).Error; err != nil {
		return nil, err
	}

	return
}

func GetIngWorld(id string) (worldList []*World, err error) {
	if err = dao.DB.Where("wstatus = ? AND userid = ?", "绘制中", id).Order("wid DESC").Find(&worldList).Error; err != nil {
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

func ConfirmPicture(background string, wid int, wsize int) (err error) {
	var world World
	fmt.Println("wid: ", wid)
	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	world.WPicture = json.RawMessage(background)
	world.WSize = wsize

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

func ConfirmEmpty(wid int, size int) (err error) {
	var world World

	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}
	world.WSize = size
	if err = dao.DB.Save(&world).Error; err != nil {
		return err
	}

	return nil
}

func DeleteWorld(wid int) error {
	var world World
	if err := dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	// 删除世界
	if err := dao.DB.Delete(&world).Error; err != nil {
		return err
	}

	return nil
}

func ChangeWorldStatus(wid int) error {
	var world World
	if err := dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	// 切换状态
	if world.WStatus == "绘制中" {
		world.WStatus = "绘制完毕"
	} else if world.WStatus == "绘制完毕" {
		world.WStatus = "绘制中"
	} else {
		return errors.New("无效的 wstatus")
	}

	// 更新 wstatus
	if err := dao.DB.Model(&world).Update("wstatus", world.WStatus).Error; err != nil {
		return err
	}

	return nil
}

func GetSearchIngWorldList(id string, key string) (worldList []*World, err error) {
	if err := dao.DB.Where("userid = ?", id).Where("wstatus = ?", "绘制中").Where("wdesc LIKE ? OR wname LIKE ?", "%"+key+"%", "%"+key+"%").Find(&worldList).Error; err != nil {
		return nil, err
	}
	return
}

func GetSearchEdWorldList(id string, key string) (worldList []*World, err error) {
	if err := dao.DB.Where("userid = ?", id).Where("wstatus = ?", "绘制完毕").Where("wdesc LIKE ? OR wname LIKE ?", "%"+key+"%", "%"+key+"%").Find(&worldList).Error; err != nil {
		return nil, err
	}
	return
}

func ModifyWorldName(name string, wid int) (err error) {
	var world World
	fmt.Println("wid: ", wid)
	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	world.WName = name

	if err = dao.DB.Save(&world).Error; err != nil {
		return err
	}

	return nil
}

func ModifyWorldDesc(desc string, wid int) (err error) {
	var world World
	fmt.Println("wid: ", wid)
	if err = dao.DB.Where("wid = ?", wid).First(&world).Error; err != nil {
		return err
	}

	world.WDesc = desc

	if err = dao.DB.Save(&world).Error; err != nil {
		return err
	}

	return nil
}
